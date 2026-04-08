# gomponents

## Objective

gomponents is an HTML5 component library written in pure Go. It lets developers build reusable, type-safe HTML components as Go functions instead of using template languages. The primary audience is Go developers building server-side rendered web applications who want compile-time guarantees, IDE support, and standard Go tooling (debugging, formatting, testing) for their HTML generation.

The library is mature and feature-complete. The core API is stable with no breaking changes. New features to the core are unlikely to be merged; the focus is on correctness, performance, and completeness of HTML5 element/attribute coverage.

## Tech Stack

- **Language**: Go (minimum version 1.18 for generics support in `Map`)
- **Dependencies**: Zero external dependencies in the core library
- **Module path**: `maragu.dev/gomponents`
- **CI**: GitHub Actions, testing against Go 1.18 through 1.25
- **Linting**: golangci-lint
- **Coverage**: Codecov, targeting 100% test coverage
- **Dependency updates**: Dependabot for GitHub Actions

## Project Structure

```
gomponents.go             # Core: Node interface, El, Attr, Text, Raw, Map, Group, If, Iff
gomponents_test.go        # Core tests and benchmarks
html/
  elements.go             # All HTML5 elements as Go functions
  attributes.go           # All HTML5 attributes as Go functions
  html_test.go            # Tests for elements and attributes
components/
  components.go           # Higher-level: HTML5 document template, Classes, JoinAttrs
  components_test.go      # Tests for components
http/
  handler.go              # HTTP handler adapter: Handler type, Adapt function
  handler_test.go         # Tests for HTTP adapter
internal/examples/app/    # Example application showing usage patterns
docs/
  spec.md                 # This file
LLMs.md                   # Documentation for LLM consumption
CLAUDE.md -> AGENTS.md    # AI assistant instructions (symlink)
CONTRIBUTING.md           # Contribution guidelines
Makefile                  # Build, test, lint, benchmark commands
.github/workflows/ci.yml  # CI configuration
.golangci.yml             # Linter configuration
```

## Commands

```sh
# Run tests with coverage and shuffled order
make test
# equivalent: go test -coverprofile=cover.out -shuffle on ./...

# Run linter
make lint
# equivalent: golangci-lint run

# Run benchmarks
make benchmark
# equivalent: go test -bench . -benchmem ./...

# View coverage report in browser
make cover
# equivalent: go tool cover -html=cover.out

# Build all packages
go build ./...
```

## Code Style

- Standard Go conventions, enforced by golangci-lint.
- Dot imports of `maragu.dev/gomponents/html` are idiomatic for this project and whitelisted in the linter config.
- HTML element and attribute function names match their HTML equivalents exactly (e.g., `Div`, `Class`, `Href`).
- Name clashes between elements and attributes are resolved by suffixing the less common usage: `Style` (attribute) / `StyleEl` (element), `Form` (element) / `FormAttr` (attribute), etc.
- Void elements (br, img, input, etc.) silently ignore non-attribute children to produce valid HTML.
- All types that render implement `fmt.Stringer` for debugging convenience.
- Render directly to `io.Writer`; use `io.StringWriter` optimization when available.
- No reflection in hot paths. No unnecessary allocations.

## Testing

- Framework: Go's standard `testing` package.
- Tests live alongside their source files (`*_test.go`).
- Run with: `go test -shuffle on ./...`
- 100% test coverage is required and enforced via Codecov.
- Table-driven tests where appropriate.
- Test both successful rendering and error cases.
- Benchmarks exist for core rendering paths and large document generation.

## Git Workflow

- Main branch: `main`.
- PRs target `main`; CI runs on push to `main` and on PRs.
- CI tests across Go 1.18 through 1.25.
- Concurrency control: CI cancels in-progress runs for the same branch.
- Contributions: small/isolated changes can go directly as PRs; larger changes should start with an issue.

## Boundaries

- **No new core API functions.** The core (`gomponents.go`) is feature-complete. Collection helpers (beyond `Map`) and flow control variants (`IfElse`, `Else`) will not be added. Users should use Go stdlib or their own helpers.
- **No external dependencies** in the core library or any sub-package.
- **No breaking changes** to the public API.
- **No template languages or code generation.** Everything is plain Go functions.
- **`logo.png`**: Do not modify or replace.
- **`LICENSE`**: Do not modify.
