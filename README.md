# Valec

[![Build Status](https://travis-ci.org/dtan4/valec.svg?branch=master)](https://travis-ci.org/dtan4/valec)
[![codecov](https://codecov.io/gh/dtan4/valec/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/valec)
[![GitHub release](https://img.shields.io/github/release/dtan4/valec.svg)](https://github.com/dtan4/valec/releases)

Handle application secrets securely

Valec is a CLI tool to handle application secrets securely using AWS DynamoDB and KMS.
Valec enables you to manage application secrets in your favorite VCS.

## Usage

### `valec dump`

Dump secrets in dotenv format

```bash
$ valec dump hoge
HOGE=fuga
```

With `-t TEMPLATE` flag, Valec dumps secrets as the form of embedding them in the given dotenv file. To override all values written in dotenv file, please specify `--override` flag too.

```bash
$ cat .env
FOO=
HOGE=hogehoge
hogehoge
YEAR=2015

# hogefugapiyo
SSSS=

$ valec dump hoge -t .env
FOO=barbarbar
HOGE=hogehoge
hogehoge
YEAR=2015

# hogefugapiyo
SSSS=
```

### `valec encrypt`

Encrypt secret

With `--add FILE` flag, encrypted secret will be added to the specified file.

```bash
$ valec encrypt NAME=awesome
AQECAHi1osu8IsEnPMo1...

$ valec encrypt NAME=awesome --add secrets.yml
$ cat secrets.yml
- key: NAME
  value: AQECAHi1osu8IsEnPMo1...
```

### `valec exec`

Execute commands using stored secrets

```bash
$ env | grep HOGE

$ valec exec hoge env | grep HOGE
HOGE=fuga
```

### `valec init`

Initialize Valec environment

These resources will be created:

- KMS key and alias (default: `valec`)
- DynamoDB table (default: `valec`)

```bash
$ valec init
```

### `valec list`

List stored secrets

```bash
# List secrets stored in DynamoDB
$ valec list hoge
HOGE: fuga

# List secrets stored in local file
$ valec list -f hoge.yaml
```

### `valec namespaces`

List all namespaces

```bash
$ valec namespaces
hoge
```

### `valec sync`

Synchronize secrets between local file and DynamoDB

```bash
$ valec sync hoge.yaml
No config will be deleted.

1 configs of hoge namespace will be added.
- PPAP
1 configs of hoge namespace was successfully added.
```

If `--dry-run` flag is given, Valec does not modify DynamoDB table actually. This might be useful for CI use.

```bash
$ valec sync hoge.yaml --dry-run
No config will be deleted.

1 configs of hoge namespace will be added.
- PPAP
```

## Development

Retrieve this repository and build using `make`.

```bash
$ go get -d github.com/dtan4/valec
$ cd $GOPATH/src/github.com/dtan4/valec
$ make deps
$ make
```

## Author

Daisuke Fujita (@dtan4)

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
