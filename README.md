# Valec

[![Build Status](https://travis-ci.org/dtan4/valec.svg?branch=master)](https://travis-ci.org/dtan4/valec)
[![Build status](https://ci.appveyor.com/api/projects/status/5re33kx1xwqgeswi/branch/master?svg=true)](https://ci.appveyor.com/project/dtan4/valec/branch/master)
[![codecov](https://codecov.io/gh/dtan4/valec/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/valec)
[![GitHub release](https://img.shields.io/github/release/dtan4/valec.svg)](https://github.com/dtan4/valec/releases)

Handle application secrets securely

Valec is a CLI tool to handle application secrets securely using AWS DynamoDB and KMS.
Valec enables you to manage application secrets in your favorite VCS.

## Workflow

1. Set up DynamoDB and KMS (first time only).

    ```bash
    $ valec init
    ```

2. Store secrets in local file. Values are encrypted.

    ```bash
    $ valec encrypt AWS_ACCESS_KEY_ID=AKIAxxxx --add production.yaml
    $ valec encrypt AWS_SECRET_ACCESS_KEY=yyyyyyyy --add production.yaml
    $ cat production.yaml
    - key: AWS_SECRET_ACCESS_KEY
      value: AQECAHi1osu...
    - key: AWS_ACCESS_KEY_ID
      value: AQECAHi1osu...
    ```

3. Save secrets to DynamoDB table.

    ```bash
    $ valec sync production.yaml
    No config will be deleted.

    2 configs of production namespace will be added.
    - AWS_SECRET_ACCESS_KEY
    - AWS_ACCESS_KEY_ID
    2 configs of production namespace were successfully added.
    ```

4. Use stored secrets in your application.

    Use stored secrets directly:

    ```bash
    $ valec exec bin/server
    ```

    or use as dotenv:

    ```bash
    $ valec dump production > .env
    $ bin/server
    ```

## Usage

### `valec dump`

Dump secrets in dotenv format

```bash
$ valec dump hoge
HOGE=fuga
```

With `-t TEMPLATE` flag, Valec dumps secrets as the form of embedding them in the given dotenv file. To override all values written in dotenv file, please specify `--override` flag too.

```bash
$ cat .env.sample
FOO=
HOGE=hogehoge
hogehoge
YEAR=2015

# comment
SSSS=

$ valec dump hoge -t .env.sample
FOO=barbarbar
HOGE=hogehoge
hogehoge
YEAR=2015

# comment
SSSS=

$ valec dump hoge -t .env.sample > .env
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

Argument must be a directory that contains secret files. `hoge.yaml` will be synchronized to `hoge` namespace.

```bash
$ ls secrets
fuga.yaml       hoge.yaml

$ valec sync secrets
fuga
  No config will be deleted.
  No config will be added.
hoge
  No config will be deleted.
  1 configs of unko namespace will be added.
    + HOGE
  1 configs of unko namespace were successfully added.
```

If `--dry-run` flag is given, Valec does not modify DynamoDB table actually. This might be useful for CI use.

```bash
$ valec sync configs --dry-run
fuga
  No config will be deleted.
  No config will be added.
hoge
  No config will be deleted.
  1 configs of unko namespace will be added.
    + HOGE
```

### `valec validate`

Validate secrets in local files

```bash
$ valec validate secrets
secrets/fuga.yaml
secrets/hoge.yaml
All configs are valid.
```

When invalid values exist:

```bash
$ valec validate secrets
secrets/fuga.yaml
secrets/hoge.yaml
  Config value is invalid. Please try `valec encrypt`. key=HOGE
Failed to validate configs. filename=tmp/unko.yaml: Some configs are invalid.
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
