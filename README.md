# gomponents

[![GoDoc](https://godoc.org/github.com/maragudk/gomponents?status.svg)](https://godoc.org/github.com/maragudk/gomponents)
[![codecov](https://codecov.io/gh/maragudk/gomponents/branch/master/graph/badge.svg)](https://codecov.io/gh/maragudk/gomponents)

gomponents are declarative view components in Go, that can render to HTML5.
gomponents aims to make it easy to build HTML5 pages of reusable components,
without the use of a template language. Think server-side-rendered React,
but without the virtual DOM and diffing.

The implementation is very usable, but the API may change until version 1 is reached.

Check out the blog post [gomponents: declarative view components in Go](https://www.maragu.dk/blog/gomponents-declarative-view-components-in-go/)
for background.

## Features

- Build reusable view components
- Write declarative HTML5 in Go without all the strings, so you get
  - Type safety
  - Auto-completion
  - Nice formatting with `gofmt`
- Simple API that's easy to learn and use (you know most already if you know HTML)
- No external dependencies

## Usage

Get the library using `go get`:

```shell script
go get -u github.com/maragudk/gomponents
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
				TitleEl(title),
				StyleEl(Type("text/css"), g.Raw(".is-active{ font-weight: bold }")),
			),
			Body(
				Navbar(currentPath),
				H1(title),
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
	return A(href, c.Classes{"is-active": currentPath == href}, g.Text(name))
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
			H1(title),
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
	return A(href, c.Classes{"is-active": currentPath == href}, g.Text(name))
}
```

For more complete examples, see [the examples directory](examples/).

### What's up with the specially named elements and attributes?

Unfortunately, there are three main name clashes in HTML elements and attributes, so they need an `El` or `Attr` suffix,
respectively, to be able to co-exist in the same package in Go:

- `form` (`FormEl`/`FormAttr`)
- `style` (`StyleEl`/`StyleAttr`)
- `title` (`TitleEl`/`TitleAttr`)
