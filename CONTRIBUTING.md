# Contributing

First of all, thank you for considering a contribution to this project. 😊
Whether it's a small fix, an idea, a bug report, or some feature code,
actions that impact the project positively are always very welcome.

## How to contribute

Feel free to submit a PR directly, if your change is in code and small and isolated.
If you're in doubt about whether your contribution is a good idea for the project,
feel free to create an issue first discussing the change.
This also applies for any larger changes; start with an issue instead of risking
a large PR that doesn't get accepted, which would make everyone involved sad.

## Development

Run the tests with `make test` and the linter with `make lint`.
Linting needs [golangci-lint](https://golangci-lint.run) v2 installed.

The library keeps 100% coverage outside `internal/`, so cover new code with tests.
Run `make cover` after `make test` to open the coverage report in a browser.

Run `make benchmark` to compare rendering performance before and after a change.
`make fuzz` runs each fuzz target for 10 seconds.

CI tests every Go version from 1.18, so avoid newer language and library features.

## Terms

By contributing code, you declare that you have the rights to add it,
and you accept that it will be published in the project under the existing license.
