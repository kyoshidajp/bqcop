# :policeman: bqcop

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: https://github.com/kyoshidajp/bqcop/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/kyoshidajp/bqcop

**bqcop** is CLI to fetch BigQuery jobs list and store it to DB.

## Usage

```
bqcop -project-id=project-id -auth-json=auth-json [options...]
```

Default) `sqlite.db` will be generated in your current directory.

### Options

```
-project-id      project id of BigQuery.

-auth-json       auth file of BigQuery.

-d, --debug      Enable debug mode.

-v, --version    Print current version.
```

### go get

If you are a Golang developper/user; then execute `go get`.

```
$ go get -u github.com/kyoshidajp/bqcop
```

## Author

[Katsuhiko YOSHIDA](https://github.com/kyoshidajp)
