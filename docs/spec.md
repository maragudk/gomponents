# gomponents

## Objective

gomponents is an HTML5 component library written in pure Go. It lets developers build reusable, type-safe HTML components using regular Go functions instead of template languages. The target audience is Go developers building server-side rendered web applications who want compile-time safety, IDE support, and standard Go tooling (debugger, `gofmt`, `goimports`) for their HTML layer.

The library is mature and stable (v1.2.0+). The API is considered feature-complete. New core features are unlikely to be accepted; the focus is on correctness, performance, and maintaining backwards compatibility.

## Tech Stack

- **Language**: Go, minimum version 1.18 (for generics in `Map`).
- **Dependencies**: Zero external dependencies in the core library.
- **CI**: GitHub Actions, testing against Go 1.18 through 1.26.
- **Linting**: golangci-lint with staticcheck.
- **Coverage**: 100% test coverage, tracked via Codecov.

## Project Structure

```
gomponents.go              Core: Node interface, El, Attr, Text, Raw, Map, Group, If, Iff
gomponents_test.go         Core tests
gomponents_benchmark_test.go  Benchmarks
html/
  elements.go             All HTML5 elements as Go functions
  attributes.go           All HTML5 attributes as Go functions
components/
  components.go           Higher-level: HTML5 document template, Classes helper, JoinAttrs
http/
  handler.go              HTTP handler adapter (Handler type, Adapt function)
x/
  slices/                 Experimental: generic Map, Filter, Reduce (no stability guarantees)
internal/
  examples/               Example application
docs/
  spec.md                 This file
  diary/                  Implementation diary
```

## Commands

```shell
# Run all tests with shuffled order and coverage
make test
# Or directly:
go test -coverprofile=cover.out -shuffle on ./...

# Run linter
make lint
# Or directly:
golangci-lint run

# Run benchmarks
make benchmark
# Or directly:
go test -bench . -benchmem ./...

# Run fuzz tests (10s per target)
make fuzz

# View coverage report in browser
make cover
```

## Code Style

- Standard Go conventions. No special formatting beyond `gofmt`.
- Tabs for indentation in Go and Markdown files (see `.editorconfig`).
- HTML element and attribute function names match their HTML equivalents exactly.
- Name clashes between elements and attributes are resolved with `El`/`Attr` suffixes. The more common usage gets the plain name (e.g., `Style` is the attribute, `StyleEl` is the element).
- Dot-imports of `gomponents/html` are permitted and whitelisted in the linter config.
- No external dependencies in the core library. Ever.
- Render directly to `io.Writer` without intermediate allocations. Use `io.StringWriter` optimization when available.
- Compile-time interface satisfaction checks for key types.

## Testing

- Framework: standard `testing` package.
- All tests live in `*_test.go` files alongside the code they test.
- Table-driven tests where appropriate.
- 100% test coverage is required and enforced via CI.
- Fuzz tests exist for complex string-processing functions (e.g., `JoinAttrs`).
- Benchmarks use `testing.B.Loop` on Go 1.24+ and include a realistic full-page rendering benchmark.
- Test both successful rendering and error cases.

## Git Workflow

- Main branch: `main`.
- CI runs on push to `main` and on pull requests targeting `main`.
- PRs are the standard contribution path. Small, isolated changes can go directly as PRs; larger changes should start with an issue.
- Commit messages are concise and descriptive, focusing on what changed and why.

## Boundaries

- **No new core flow-control functions.** `IfElse`, `Else`, and similar will not be added. Users should use Go's own control flow or IIFEs.
- **No breaking changes.** The library is stable. Public API changes must be backwards-compatible.
- **No external dependencies** in the core library (`gomponents`, `html`, `components`, `http` packages).
- **The `x/` packages** are experimental and do not carry the same stability guarantees.
- **Do not modify** the logo (`logo.png`), license (`LICENSE`), or CI secrets configuration.
