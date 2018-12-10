# :policeman: bqcop

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: https://github.com/kyoshidajp/bqcop/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/kyoshidajp/bqcop

**bqcop** is CLI to fetch BigQuery jobs and store it to DB.

## Usage

```
bqcop -project-id=project-id -auth-json=auth-json [options...]
```

Fetch BigQuery jobs executed during the **24 hours from now** by calling [Jobs list API](https://cloud.google.com/bigquery/docs/reference/rest/v2/jobs/list) and store it to DB.

### Options

```
-project-id      Project ID of BigQuery.

-auth-json       Auth File of BigQuery.

-db-dialect      Dialect of Database.
                 default: sqlite3

-db-path         Path of Database.
                 default: sqlite.db

-d, --debug      Enable debug mode.

-v, --version    Print current version.
```

### Output

`sqlite.db` which has `bq_jobs` will be generated in your current directory if both `-db-dialect` and `-db-path` are not specified.

Schema of `bq_jobs` table.

| field | type | description |
| ----- | ---- | --- |
| id | integer | primary key |
| created_at | datetime | created time |
| updated_at | datetime | updated time |
| deleted_at | datetime | deleted time |
| job_id | varchar(255) | job id |
| query | varchar(255) | job query |
| user_email | varchar(255) | user who exec query |
| total_bytes_billed | bigint | total bytes of billed |
| start_time | datetime | job started time |
| end_time | datetime | job ended time |

## Install

### go get

If you are a Golang developper/user; then execute `go get`.

```
$ go get -u github.com/kyoshidajp/bqcop
```

## Author

[Katsuhiko YOSHIDA](https://github.com/kyoshidajp)
