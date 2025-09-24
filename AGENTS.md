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
