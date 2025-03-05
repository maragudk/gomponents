# Tired of complex template languages?

<img src="logo.png" alt="Logo" width="300" align="right">

[![GoDoc](https://pkg.go.dev/badge/maragu.dev/gomponents)](https://pkg.go.dev/maragu.dev/gomponents)
[![Go](https://github.com/maragudk/gomponents/actions/workflows/ci.yml/badge.svg)](https://github.com/maragudk/gomponents/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/maragudk/gomponents/branch/main/graph/badge.svg)](https://codecov.io/gh/maragudk/gomponents)
[![Go Report Card](https://goreportcard.com/badge/maragu.dev/gomponents)](https://goreportcard.com/report/maragu.dev/gomponents)

Try HTML components in pure Go.

_gomponents_ are HTML components written in pure Go.
They render to HTML 5, and make it easy for you to build reusable components.
So you can focus on building your app instead of learning yet another templating language.

```shell
go get maragu.dev/gomponents
```

Made with ✨sparkles✨ by [maragu](https://www.maragu.dev/).

Does your company depend on this project? [Contact me at markus@maragu.dk](mailto:markus@maragu.dk?Subject=Supporting%20your%20project) to discuss options for a one-time or recurring invoice to ensure its continued thriving.

## Features

Check out [www.gomponents.com](https://www.gomponents.com) for an introduction.

- Build reusable HTML components
- Write declarative HTML 5 in Go without all the strings, so you get
  - Type safety from the compiler
  - Auto-completion from the IDE
  - Easy debugging with the standard Go debugger
  - Automatic formatting with `gofmt`/`goimports`
- Simple API that's easy to learn and use (you know most already if you know HTML)
- Useful helpers like
  - `Text` and `Textf` that insert HTML-escaped text,
  - `Raw` and `Rawf` for inserting raw strings,
  - `Map` for mapping data to components and `Group` for grouping components,
  - and `If`/`Iff` for conditional rendering.
- No external dependencies
- Mature and stable, no breaking changes

## Usage

```shell
go get maragu.dev/gomponents
```

```go
package main

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Navbar(authenticated bool, currentPath string) Node {
	return Nav(
		NavbarLink("/", "Home", currentPath),
		NavbarLink("/about", "About", currentPath),
		If(authenticated, NavbarLink("/profile", "Profile", currentPath)),
	)
}

func NavbarLink(href, name, currentPath string) Node {
	return A(Href(href), Classes{"is-active": currentPath == href}, g.Text(name))
}
```

(Some people don't like dot-imports, and luckily it's completely optional.)

For a more complete example, see [the examples directory](internal/examples/).
There's also the [gomponents-starter-kit](https://github.com/maragudk/gomponents-starter-kit) for a full application template.

## Architecture

gomponents is organized into several packages:

- `gomponents`: Core interfaces and functions like `Node`, `El`, `Attr`, and helpers like `Map`, `Group`, `If`, `Text`, `Raw`.
- `gomponents/html`: HTML elements and attributes.
- `gomponents/components`: Higher-level components and utilities.
- `gomponents/http`: HTTP-related utilities for web servers.

### Void Elements

Void elements in HTML (like `<br>`, `<img>`, `<input>`) don't have closing tags.
gomponents handles these correctly by checking against an internal list of void elements during rendering.
When you create a void element, any child nodes that are not attributes will be ignored automatically to ensure valid HTML output.

## Performance Considerations

gomponents renders directly to an `io.Writer`, making it efficient for server-side rendering.
The library avoids unnecessary allocations where possible.

## FAQ

### Is gomponents production-ready?

Yes! gomponents is mature, stable, fully tested with 100% coverage, and has been used in production by myself and many others.

### Should I choose `html/template`, Templ, or gomponents?

These are all good choices, and it largely comes down to preference.
I wrote gomponents because I didn't like how I think it's hard to pass data around between templates in `html/template`.
gomponents is pure Go, with no extra build step like Templ, so it works with all tools that already support Go.

That said, both `html/template` and Templ will do the same thing as gomponents in the end. Try them all and choose what you like!

### I don't like how HTML looks in Go.

First of all, that's not a question. 😉

More seriously, think of gomponents like a DSL for HTML. You're building UI components. Give it a day, and it'll feel natural.

### What's up with the specially named elements and attributes?

Unfortunately, there are some name clashes in HTML elements and attributes, so they need an `El` or `Attr` suffix,
to be able to co-exist in the same package in Go.

I've chosen one or the other based on what I think is the common usage.
In either case, the less-used variant also exists in the codebase:

- `cite` (`Cite`/`CiteAttr`, `CiteEl` also exists)
- `data` (`DataEl`/`Data`, `DataAttr` also exists)
- `form` (`Form`/`FormAttr`, `FormEl` also exists)
- `label` (`Label`/`LabelAttr`, `LabelEl` also exists)
- `style` (`StyleEl`/`Style`, `StyleAttr` also exists)
- `title` (`TitleEl`/`Title`, `TitleAttr` also exists)
