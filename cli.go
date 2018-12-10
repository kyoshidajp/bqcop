package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/colorstring"
	"google.golang.org/api/option"
)

const (
	gcpKey string = "auth.json"
	dbFile string = "sqlite.db"
)

const (
	// EnvDebug is environmental var to handle debug mode
	EnvDebug = "BQCOP_DEBUG"
)

// Exit codes are in value that represnet an exit code for a paticular error
const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeBadArgs
)

// CLI is the command line object
type CLI struct {
	outStream, errStream io.Writer
}

// BQJob is BigQuery Job Object
type BQJob struct {
	gorm.Model
	JobID            string `gorm:"primary_key"`
	Query            string
	UserEmail        string
	TotalBytesBilled int64
	StartTime        time.Time
	EndTime          time.Time
}

// BQCop is exec client
type BQCop struct {
	db        *gorm.DB
	authFile  string
	projectID string
}

// Debugf prints debug output when EnvDebug is given
func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// PrintErrorf prints error message on console
func PrintErrorf(format string, args ...interface{}) {
	format = fmt.Sprintf("[red]%s[reset]\n", format)
	fmt.Fprint(colorable.NewColorableStderr(),
		colorstring.Color(fmt.Sprintf(format, args...)))
}

// insert does
func (b *BQCop) insert(ctx context.Context, it *bigquery.JobIterator) *bigquery.JobIterator {
	job, err := it.Next()
	if err != nil {
		return nil
	}
	if job == nil {
		return nil
	}

	status, err := job.Status(ctx)
	if err != nil || status == nil {
		return it
	}

	jc, _ := job.Config()

	if jc == nil {
		return it
	}

	queryConfig, ok := jc.(*bigquery.QueryConfig)
	if !ok {
		return it
	}

	query := queryConfig.Q
	if query == "" {
		return it
	}

	id := job.ID()
	stat := status.Statistics
	bqJob := &BQJob{
		JobID:            id,
		Query:            query,
		UserEmail:        job.Email(),
		TotalBytesBilled: stat.TotalBytesProcessed,
		StartTime:        stat.StartTime,
		EndTime:          stat.EndTime,
	}

	log.Printf("Insert jobID: %s", id)

	b.db.Create(bqJob)
	return it
}

// Do does do
func (b *BQCop) Do() int {
	ctx := context.Background()
	client, err := bigquery.NewClient(
		ctx,
		b.projectID,
		option.WithServiceAccountFile(b.authFile))
	if err != nil {
		panic("ERROR")
	}
	defer client.Close()

	it := client.Jobs(ctx)
	it.State = bigquery.Done
	it.AllUsers = true

	now := time.Now()
	it.MinCreationTime = now.AddDate(0, 0, -1) // 1 day ago
	it.MaxCreationTime = now

	for {
		it := b.insert(ctx, it)
		if it == nil {
			break
		}
	}
	log.Printf("Finished.")

	defer b.db.Close()
	return ExitCodeOK
}

// connectDB connects data base
func connectDB(dialect string, path string) *gorm.DB {
	Debugf("Connect to: %s %s", dialect, path)

	db, err := gorm.Open(dialect, path)
	if err != nil {
		// done
		return nil
	}
	db.AutoMigrate(&BQJob{})

	return db
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	var (
		auth      string
		projectID string
		dbDialect string
		dbPath    string
		debug     bool
		version   bool
	)
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}
	flags.StringVar(&auth, "auth-json", "", "")
	flags.StringVar(&projectID, "project-id", "", "")
	flags.StringVar(&dbDialect, "db-dialect", "sqlite3", "")
	flags.StringVar(&dbPath, "db-path", dbFile, "")
	flags.BoolVar(&debug, "debug", false, "")
	flags.BoolVar(&debug, "d", false, "")
	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	// Parse flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if debug {
		os.Setenv(EnvDebug, "1")
		Debugf("Run as DEBUG mode")
	}

	if version {
		fmt.Fprintf(c.outStream, fmt.Sprintf("%s\n", Version))
		return ExitCodeOK
	}

	authFile := gcpKey
	if auth != "" {
		authFile = auth
	}

	if projectID == "" {
		PrintErrorf("Invalid argument: -project-id must be specified.")
		return ExitCodeBadArgs
	}

	db := connectDB(dbDialect, dbPath)
	if db == nil {
		return ExitCodeError
	}

	cop := &BQCop{
		db:        db,
		authFile:  authFile,
		projectID: projectID,
	}
	return cop.Do()
}

var helpText = `Usage: bqcop -project-id=project-id -auth-json=auth-json [options...]

bqcop is a tool to fetch BigQuery jobs and store it to DB.

Options:

  -project-id      project id of BigQuery.

  -auth-json       auth file of BigQuery.

  -d, --debug      Enable debug mode.

  -v, --version    Print current version.
`
