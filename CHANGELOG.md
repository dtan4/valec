# [v0.5.3](https://github.com/dtan4/valec/releases/tag/v0.5.3) (2017-01-17)

## Fixes

- Divide BatchWriteItem API requests for many items [#44](https://github.com/dtan4/valec/pull/46)

# [v0.5.2](https://github.com/dtan4/valec/releases/tag/v0.5.2) (2016-12-27)

## Fixes

- Delete namespace correctly [#44](https://github.com/dtan4/valec/pull/44)

# [v0.5.1](https://github.com/dtan4/valec/releases/tag/v0.5.1) (2016-12-26)

## Features

- Delete namespace which is no longer managed [#41](https://github.com/dtan4/valec/pull/41)
- Calculate namespace from full path [#40](https://github.com/dtan4/valec/pull/40)

# [v0.5.0](https://github.com/dtan4/valec/releases/tag/v0.5.0) (2016-12-22)

## Backward incompatible changes

- Secret YAML file schema was changed [#38](https://github.com/dtan4/valec/pull/38)

## Features

- Specify KMS key per namespace [#38](https://github.com/dtan4/valec/pull/38)
- Interactive encryption mode (`valec encrypt -i KEY1 KEY2`) [#34](https://github.com/dtan4/valec/pull/34)
- Use aws-sdk-go 1.6.x [#33](https://github.com/dtan4/valec/pull/33)
- Encrypt secrets from stdin [#32](https://github.com/dtan4/valec/pull/32)
- Add `valec dotenv` command [#31](https://github.com/dtan4/valec/pull/31)
- Encrypt multiple secrets at once (`valec encrypt KEY1=VALUE1 KEY2=VALUE2`) [#30](https://github.com/dtan4/valec/pull/30)

# [v0.4.0](https://github.com/dtan4/valec/releases/tag/v0.4.0) (2016-12-21)

## Features

- Interactive encryption mode (`valec encrypt -i KEY1 KEY2`) [#34](https://github.com/dtan4/valec/pull/34)
- Use aws-sdk-go 1.6.x [#33](https://github.com/dtan4/valec/pull/33)
- Encrypt secrets from stdin [#32](https://github.com/dtan4/valec/pull/32)
- Add `valec dotenv` command [#31](https://github.com/dtan4/valec/pull/31)
- Encrypt multiple secrets at once (`valec encrypt KEY1=VALUE1 KEY2=VALUE2`) [#30](https://github.com/dtan4/valec/pull/30)

# [v0.3.2](https://github.com/dtan4/valec/releases/tag/v0.3.2) (2016-12-16)

## Features

- Validate nested secret files [#28](https://github.com/dtan4/valec/pull/28)
- Add `valec dump -q` flag to dump dotenv values as quoted string [#27](https://github.com/dtan4/valec/pull/27)

# [v0.3.1](https://github.com/dtan4/valec/releases/tag/v0.3.1) (2016-12-15)

## Features

- Add `valec dump --output` flag to dump dotenv with preserved lines [#25](https://github.com/dtan4/valec/pull/25)

# [v0.3.0](https://github.com/dtan4/valec/releases/tag/v0.3.0) (2016-12-12)

## Features

- Show error if the given namespace does not exist [#22](https://github.com/dtan4/valec/pull/22)
- Detect updated changes [#21](https://github.com/dtan4/valec/pull/21)
- Sort secrets alphabetically [#20](https://github.com/dtan4/valec/pull/20)
- Synchronize nested namespace [#19](https://github.com/dtan4/valec/pull/19)

# [v0.2.1](https://github.com/dtan4/valec/releases/tag/v0.2.1) (2016-12-05)

## Features

- Add `--region` flag to specify AWS region [#16](https://github.com/dtan4/valec/pull/16), [#17](https://github.com/dtan4/valec/pull/17)

# [v0.2.0](https://github.com/dtan4/valec/releases/tag/v0.2.0) (2016-12-02)

## Backward incompatible changes

- Synchronize local files in the given directory by `valic sync` [#13](https://github.com/dtan4/valec/pull/13)
  - `--namespace` option was deprecated.
  - Argument was changed from a file to a directory.

## Features

- Add new command `valec validate` [#14](https://github.com/dtan4/valec/pull/14)

# [v0.1.1](https://github.com/dtan4/valec/releases/tag/v0.1.1) (2016-12-01)

## Features

- Add dry-run feature to `valec sync` command [#10](https://github.com/dtan4/valec/pull/10)

# [v0.1.0](https://github.com/dtan4/valec/releases/tag/v0.1.0) (2016-11-30)

Initial release.
