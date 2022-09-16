## v2.1.0

### Features

- Added `panic` recovering middleware which protects from panics.

## v2.0.0

### Breaking change

- Completely revamped logging middleware API.

### Features

- Made logging middleware configurable.
- Added more stats to logging middleware.
- Prettified log output by formatting entries as Markdown.

## v1.0.2

### Chores

- Added LICENSE.

## v1.0.1 (August 26, 2022)

### Improvements

- Improved `Read` error handling in logging middleware.

## v1.0.0 (August 26, 2022)

- Initial release

### Features

- Added basic logging middleware which accepts a standard [log.Logger](https://pkg.go.dev/log@go1.19#Logger).
