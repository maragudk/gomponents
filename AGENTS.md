# AI development guide

This is a guide for how AIs should develop code for Markus.

## About me

I'm Markus. Call me that.
I'm an independent software consultant, developing software and digital products professionally for 10+ years.
I specialize in cloud-native Go application development as well as AI engineering.

Some highlights about my professional perspective on my work:
- I work almost exclusively with Go. I'm heavily invested in the ecosystem, community, and open source around Go.
- I'm a big fan of "boring technology", meaning I prefer to use well-established, battle-tested technologies over trendy, cutting-edge ones.
- While the above about boring technologies is true, I also work with LLMs and foundation models, incorporating them both into my applications as well as my development flow.
- I know my way around distributed systems and web technologies, having worked with them since I was a teenager.
- I prefer SQLite and PostgreSQL for databases. I also like object stores (such as S3), queues, and load balancers, but generally don't use any other cloud primitives.
- I don't like microservice-oriented architectures, preferring a more monolithic approach.

## Development style

### Go application structure

Generally, I build web applications and libraries/modules.

These are the packages typically present in applications (some may be missing, which typically means I don't need them in the project).

- `main`: contains the main entry point of the application (in directory `cmd/app`)
- `model`: contains the domain model used throughout the other packages
- `sql`/`sqlite`/`postgres`: contains SQL database-related logic as well as database migrations (under subdirectory `migrations/`). The database used is either SQLite or PostgreSQL.
- `sqltest`/`sqlitetest`/`postgrestest`: package used in testing, for setting up and tearing down test databases
- `s3`: logic for interacting with Amazon S3 or compatible object stores
- `s3test`: package used in testing, for setting up and tearing down test S3 buckets
- `llm`: clients for interacting with large language models (LLMs) and foundation models
- `llmtest`: package used in testing, for setting up LLM clients for testing
- `http`: HTTP handlers for the application
- `html`: HTML templates for the application, written with the gomponents library (see https://www.gomponents.com/llms.txt for how to use that if you need to)

### Code style

#### Dependency injection

I make heavy use of dependency injection between components. This is typically done with private interfaces on the receiving side. Note the use of `userGetter` in this example:

```go user.go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"maragu.dev/httph"

	"model"
)

type UserResponse struct {
	Name string
}

type userGetter interface {
	GetUser(ctx context.Context, id model.ID) (model.User, error)
}

func User(r chi.Router, db userGetter) {
	r.Get("/user", httph.JSONHandler(func(w http.ResponseWriter, r *http.Request, _ any) (UserResponse, error) {
		id := r.URL.Query().Get("id")
		user, err := db.GetUser(r.Context(), model.ID(id))
		if err != nil {
			return UserResponse{}, httph.HTTPError{Code: http.StatusInternalServerError, Err: errors.New("error getting user")}
		}
		return UserResponse{Name: user.Name}, nil
	}))
}

```

#### Tests

I write tests for most functions and methods. I almost always use subtests with a good description of whats is going on and what the expected result is.

Here's an example:

```go example.go
package example

type Thing struct {}

func (t *Thing) DoSomething() (bool, error) {
	return true, nil
}
```

```go example_test.go
package example_test

import (
	"testing"

	"maragu.dev/is"

	"example"
)

func TestThing_DoSomething(t *testing.T) {
	t.Run("should do something and return a nil error", func(t *testing.T) {
		thing := &example.Thing{}

		ok, err := thing.DoSomething()
		is.NotError(t, err)
		is.True(t, ok)
	})
}
```

Sometimes I use table-driven tests:

```go example.go
package example

import "errors"

type Thing struct {}

var ErrChairNotSupported = errors.New("chairs not supported")

func (t *Thing) DoSomething(with string) error {
	if with == "chair" {
		return ErrChairNotSupported
	}
	return nil
}
```

```go example_test.go
package example_test

import (
	"testing"

	"maragu.dev/is"

	"example"
)

func TestThing_DoSomething(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{name: "should do something with the table and return a nil error", input: "table", expected: nil},
		{name: "should do something with the chair and return an ErrChairNotSupported", input: "chair", expected: example.ErrChairNotSupported},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			thing := &example.Thing{}

			err := thing.DoSomething(test.input)
			if test.expected != nil {
				is.Error(t, test.expected, err)
			} else {
				is.NotError(t, err)
			}
		})
	}
}
```

I prefer integration tests with real dependencies over mocks, because there's nothing like the real thing. Dependencies are typically run in Docker containers. You can assume the dependencies are running when running tests.

It makes sense to use mocks when the important part of a test isn't the dependency, but it plays a smaller role. But for example, when testing database methods, a real underlying database should be used.

I use test assertions with the module `maragu.dev/is`. Available functions: `is.True`, `is.Equal`, `is.Nil`, `is.NotNil`, `is.EqualSlice`, `is.NotError`, `is.Error`. All of these take an optional message as the last parameter.

Since tests are shuffled, don't rely on test order, even for subtests.

Every time the `postgrestest.NewDatabase(t)`/`sqlitetest.NewDatabase(t)` test helpers are called, the database is in a clean state (no leftovers from other tests etc.).

#### Miscellaneous

- Variable naming:
  - `req` for requests, `res` for responses
- Prefer lowercase SQL queries
- There are SQL helpers available, at `Database.H.Select`, `Database.H.Exec`, `Database.H.Get`, `Database.H.InTx`.
- Use the `any` builtin in Go instead of `interface{}`
- There's an alias for `sql.ErrNoRows` from stdlib at `maragu.dev/glue/sql.ErrNoRows`, so you don't have to import both
- In tests, use `t.Context()` instead of `context.Background()`
- Test helper functions should call `testing.T.Helper()`
- All HTML buttons need the `cursor-pointer` CSS class
- SQLite time format is always a string returned by `strftime('%Y-%m-%dT%H:%M:%fZ')`
- Remember that private functions in Go are package-level, so you can use them across files in the same package
- Lowercase the beginning of HTML component names unless they need to be used by an HTTP handler outside the package
- Documentation should follow the Go style of having the identifier name be the first word of the sentence, and then completing the sentence without repeating itself. Example: "// SearchProducts using the given search query and result limit." NOT: "// SearchProducts searches products using the given search query and result label."

### Testing, linting, evals

Run `make test` or `go test -shuffle on ./...` to run all tests. To run tests in just one package, use `go test -shuffle on ./path/to/package`. To run a specific test, use `go test ./path/to/package -run TestName`.

Run `make lint` or `golangci-lint run` to run linters. They should always be run on the package/directory level, it won't work with single files.

Run `make eval` or `go test -shuffle on -run TestEval ./...` to run LLM evals.

Run `make fmt` to format all code in the project, which is useful as a last finishing touch.

You can access the database by using `psql` or `sqlite3` in the shell.

### Version control

When writing commit messages, surround identifier names (variable names, type names, etc.) in backticks.

### Bugs

If you think you've found a bug during testing, ask me what to do, instead of trying to work around the bug in tests.

### Documentation

You can generally look up documentation for a Go module using `go doc` with the module name. For example, `go doc net/http` for something in the standard library, or `go doc maragu.dev/gai` for a third-party module. You can also look up more specific documentation for an identifier with something like `go doc maragu.dev/gai.ChatCompleter`, for the `ChatCompleter` interface.

### Checking apps in a browser

You can assume the app is running and available in a browser using the Playwright tool. It auto-reloads on code changes so you don't have to.
Log output from the running application is in `app.log` in the project root.
# gomponents Development Guide

This is gomponents, an HTML component library written in pure Go that renders to HTML5. This guide provides instructions for AI assistants working on this codebase.

## About gomponents

gomponents enables building HTML components using pure Go functions instead of template languages. Key features:
- Type-safe HTML generation with compile-time guarantees
- No external dependencies in the core library
- Direct rendering to `io.Writer` for efficiency
- Support for all HTML5 elements and attributes
- Conditional rendering and data mapping helpers

## Project Structure

The project is organized into focused packages:

- **Core (`gomponents.go`)**: Main interfaces (`Node`), element/attribute creators (`El`, `Attr`), text rendering (`Text`, `Raw`), and helpers (`Map`, `Group`, `If`, `Iff`)
- **html/**: All HTML5 elements and attributes as Go functions
- **components/**: Higher-level components like `HTML5` document structure and `Classes` helper
- **http/**: HTTP handler utilities for web servers
- **internal/examples/app/**: Example application showing usage patterns

## Development Standards

### Code Style
- Follow standard Go conventions
- Use clear, descriptive function names
- No external dependencies in core library
- Maintain backwards compatibility (library is stable/mature)
- HTML element/attribute names match their HTML equivalents exactly

### Testing
- Run tests: `make test` or `go test -shuffle on ./...`
- Run linting: `make lint` or `golangci-lint run`
- Maintain 100% test coverage
- Use table-driven tests where appropriate
- Test both successful rendering and error cases

### Performance Considerations
- Render directly to `io.Writer` without intermediate allocations
- Use `io.StringWriter` optimization when available
- Avoid reflection in hot paths
- Keep void element checks efficient

## Key Concepts

### Node Interface
Everything implements the core `Node` interface:
```go
type Node interface {
    Render(w io.Writer) error
}
```

### Node Types
- `ElementType`: HTML elements and text content
- `AttributeType`: HTML attributes (render in different phase)

### Void Elements
Self-closing HTML elements (br, img, input, etc.) are handled specially - non-attribute children are ignored to ensure valid HTML.

### Attribute vs Element Disambiguation
Some HTML names conflict (e.g., `style`). Convention:
- Most common usage gets the simple name (`Style` for attribute)
- Alternative gets suffix (`StyleEl` for element)
- Both variants always exist

## Common Patterns

### Creating Elements
```go
// Basic element
Div(Class("container"), Text("Hello"))

// Custom element
El("custom-element", Attr("data-value", "123"))
```

### Conditional Rendering
```go
If(condition, someNode)       // Eager evaluation
Iff(condition, func() Node {  // Lazy evaluation
    return expensiveNode()
})
```

### Data Mapping
```go
Map(items, func(item Item) Node {
    return Li(Text(item.Name))
})
```

## Testing Guidelines

Test both the structure and actual HTML output:
```go
func TestComponent(t *testing.T) {
    node := MyComponent("test")

    var buf bytes.Buffer
    err := node.Render(&buf)
    // Check err and buf.String()
}
```

## HTML Generation Best Practices

1. **Always escape user content**: Use `Text()` for user data, `Raw()` only for trusted HTML
2. **Leverage type safety**: Create typed component functions rather than generic ones
3. **Use Groups for multiple nodes**: Return `Group{node1, node2}` when multiple nodes needed
4. **Handle nil nodes gracefully**: The library safely ignores nil nodes during rendering

## Contributing Guidelines

- New HTML elements/attributes should follow HTML5 spec exactly
- Core library changes require careful consideration of backwards compatibility
- Performance optimizations welcome, but measure first
- Documentation should be clear and include examples
- All changes must include comprehensive tests

This is a mature, stable library focused on simplicity and performance. Prefer clear, straightforward implementations over complex features.
