# :policeman: bqcop

[![GitHub release](https://img.shields.io/github/release/kyoshidajp/bqcop.svg?style=flat-square)][release]
[![Travis](https://travis-ci.org/kyoshidajp/bqcop.svg?branch=master)](https://travis-ci.org/kyoshidajp/bqcop)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/kyoshidajp/bqcop/releases
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

### Manual

1. Download binary which meets your system from [Releases](https://github.com/kyoshidajp/bqcop/releases).
1. Unarchive it.
1. Put `bqcop` where you want.
1. Add `bqcop` path to `$PATH`.

## Author

[Katsuhiko YOSHIDA](https://github.com/kyoshidajp)
