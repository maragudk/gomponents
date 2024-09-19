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

The preferred way to use gomponents is with so-called dot-imports (note the dot before the `gomponents/html` import),
to give you that smooth, native HTML feel:

```go
package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_ = Page("Hi!", r.URL.Path).Render(w)
}

func Page(title, currentPath string) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			Head(
				TitleEl(g.Text(title)),
				StyleEl(Type("text/css"), g.Raw(".is-active{ font-weight: bold }")),
			),
			Body(
				Navbar(currentPath),
				H1(g.Text(title)),
				P(g.Textf("Welcome to the page at %v.", currentPath)),
			),
		),
	)
}

func Navbar(currentPath string) g.Node {
	return Nav(
		NavbarLink("/", "Home", currentPath),
		NavbarLink("/about", "About", currentPath),
	)
}

func NavbarLink(href, name, currentPath string) g.Node {
	return A(Href(href), c.Classes{"is-active": currentPath == href}, g.Text(name))
}
```

Some people don't like dot-imports, and luckily it's completely optional.
If you don't like dot-imports, just use regular imports.

You could also use the provided HTML5 document template to simplify your code a bit:

```go
package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_ = Page("Hi!", r.URL.Path).Render(w)
}

func Page(title, currentPath string) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			StyleEl(Type("text/css"), g.Raw(".is-active{ font-weight: bold }")),
		},
		Body: []g.Node{
			Navbar(currentPath),
			H1(g.Text(title)),
			P(g.Textf("Welcome to the page at %v.", currentPath)),
		},
	})
}

func Navbar(currentPath string) g.Node {
	return Nav(
		NavbarLink("/", "Home", currentPath),
		NavbarLink("/about", "About", currentPath),
	)
}

func NavbarLink(href, name, currentPath string) g.Node {
	return A(Href(href), c.Classes{"is-active": currentPath == href}, g.Text(name))
}
```

For more complete examples, see [the examples directory](examples/).

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
