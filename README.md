# Tired of complex template languages?

<img src="logo.png" alt="Logo" width="300" align="right"/>

[![GoDoc](https://pkg.go.dev/badge/github.com/maragudk/gomponents)](https://pkg.go.dev/github.com/maragudk/gomponents)
[![Go](https://github.com/maragudk/gomponents/actions/workflows/ci.yml/badge.svg)](https://github.com/maragudk/gomponents/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/maragudk/gomponents/branch/main/graph/badge.svg)](https://codecov.io/gh/maragudk/gomponents)
[![Go Report Card](https://goreportcard.com/badge/github.com/maragudk/gomponents)](https://goreportcard.com/report/github.com/maragudk/gomponents)

Try HTML components in pure Go.

_gomponents_ are HTML components written in pure Go.
They render to HTML 5, and make it easy for you to build reusable components.
So you can focus on building your app instead of learning yet another templating language.

The API may change until version 1 is reached.

Check out [www.gomponents.com](https://www.gomponents.com) for an introduction.

Made in 🇩🇰 by [maragu](https://www.maragu.dk), maker of [online Go courses](https://www.golang.dk/).

## Features

- Build reusable HTML components
- Write declarative HTML5 in Go without all the strings, so you get
  - Type safety
  - Auto-completion
  - Nice formatting with `gofmt`
- Simple API that's easy to learn and use (you know most already if you know HTML)
- Useful helpers like `Text` and `Textf` that insert HTML-escaped text, `Map` for mapping data to components,
  and `If`/`Iff` for conditional rendering.
- No external dependencies

## Usage

Get the library using `go get`:

```shell
go get github.com/maragudk/gomponents
```

The preferred way to use gomponents is with so-called dot-imports (note the dot before the imports),
to give you that smooth, native HTML feel:

```go
package main

import (
	. "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
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

Some people don't like dot-imports, and luckily it's completely optional.

For a more complete example, see [the examples directory](internal/examples/).

### What's up with the specially named elements and attributes?

Unfortunately, there are six main name clashes in HTML elements and attributes, so they need an `El` or `Attr` suffix,
to be able to co-exist in the same package in Go. I've chosen one or the other based on what I think is the common usage.
In either case, the less-used variant also exists in the codebase:

- `cite` (`Cite`/`CiteAttr`, `CiteEl` also exists)
- `data` (`DataEl`/`Data`, `DataAttr` also exists)
- `form` (`Form`/`FormAttr`, `FormEl` also exists)
- `label` (`Label`/`LabelAttr`, `LabelEl` also exists)
- `style` (`StyleEl`/`Style`, `StyleAttr` also exists)
- `title` (`TitleEl`/`Title`, `TitleAttr` also exists)
