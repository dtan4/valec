# Valec

[![Build Status](https://travis-ci.org/dtan4/valec.svg?branch=master)](https://travis-ci.org/dtan4/valec)

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
