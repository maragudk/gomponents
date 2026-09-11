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

- **Core (`gomponents.go`)**: Main interfaces (`Node`), element/attribute creators (`El`, `Attr`), text rendering (`Text`/`Textf`, `Raw`/`Rawf`), and helpers (`Map`, `Group`, `If`, `Iff`)
- **html/**: All HTML5 elements and attributes as Go functions
- **components/**: Higher-level components like `HTML5` document structure, `Classes` and `JoinAttrs` helpers, and more
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
- Run benchmarks: `make benchmark` or `go test -bench . -benchmem ./...`
- Run fuzzing: `make fuzz`
- Maintain 100% test coverage
- Use table-driven tests where appropriate
- Test both successful rendering and error cases

### Performance Considerations

- Render directly to `io.Writer` without intermediate allocations
- Use `io.StringWriter` optimization when available
- Avoid reflection in hot paths
- Keep void element checks efficient

## Contributing Guidelines

- New HTML elements/attributes should follow HTML5 spec exactly
- Core library changes require careful consideration of backwards compatibility
- Performance optimizations welcome, but measure first
- Documentation should be clear and include examples
- All changes must include comprehensive tests

This is a mature, stable library focused on simplicity and performance. Prefer clear, straightforward implementations over complex features.
