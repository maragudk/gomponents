# gomponents

[![GoDoc](https://godoc.org/github.com/maragudk/gomponents?status.svg)](https://godoc.org/github.com/maragudk/gomponents)
[![codecov](https://codecov.io/gh/maragudk/gomponents/branch/master/graph/badge.svg)](https://codecov.io/gh/maragudk/gomponents)

gomponents are declarative view components in Go, that can render to HTML.
gomponents aims to make it easy to build HTML pages of reusable components,
without the use of a template language. Think server-side-rendered React,
but without the virtual DOM and diffing.

The implementation is still incomplete, but usable. The API may change until version 1 is reached.

Check out the blog post [gomponents: declarative view components in Go](https://www.maragu.dk/blog/gomponents-declarative-view-components-in-go/)
for background.

## Usage

Get the library using `go get`:

```shell script
go get -u github.com/maragudk/gomponents
```

Then do something like this:

```go
package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", handler())
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page("Hi!", r.URL.Path)
		_ = g.Write(w, page)
	}
}

func Page(title, path string) g.Node {
	return el.Document(
		el.HTML(
			g.Attr("lang", "en"),
			el.Head(
				el.Title(title),
				el.Style(g.Attr("type", "text/css"), g.Raw(".is-active{font-weight: bold}")),
			),
			el.Body(
				Navbar(path),
				el.H1(title),
				el.P(g.Textf("Welcome to the page at %v.", path)),
			),
		),
	)
}

func Navbar(path string) g.Node {
	return g.El("nav",
		el.A("/", attr.Classes{"is-active": path == "/"}, g.Text("Home")),
		el.A("/about", attr.Classes{"is-active": path == "/about"}, g.Text("About")),
	)
}
```

You could also use a page template to simplify your code a bit:

```go
package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	c "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/el"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", handler())
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page("Hi!", r.URL.Path)
		_ = g.Write(w, page)
	}
}

func Page(title, path string) g.Node {
	return c.HTML5(c.DocumentProps{
		Title:       title,
		Language:    "en",
		Head:        []g.Node{el.Style(g.Attr("type", "text/css"), g.Raw(".is-active{font-weight: bold}"))},
		Body:        []g.Node{
			Navbar(path),
			el.H1(title),
			el.P(g.Textf("Welcome to the page at %v.", path)),
		},
	})
}

func Navbar(path string) g.Node {
	return g.El("nav",
		el.A("/", attr.Classes{"is-active": path == "/"}, g.Text("Home")),
		el.A("/about", attr.Classes{"is-active": path == "/about"}, g.Text("About")),
	)
}
```

For more complete examples, see [the examples directory](examples/).
