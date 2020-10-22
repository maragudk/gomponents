// Package el provides shortcuts and helpers to common HTML elements.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element for a list of elements.
package el

import (
	"strings"

	g "github.com/maragudk/gomponents"
)

// Document returns an special kind of Node that prefixes its children with the string "<!doctype html>".
func Document(children ...g.Node) g.NodeFunc {
	return func() string {
		var b strings.Builder
		b.WriteString("<!doctype html>")
		for _, c := range children {
			b.WriteString(c.Render())
		}
		return b.String()
	}
}

// HTML returns an element with name "html" and the given children.
func HTML(children ...g.Node) g.NodeFunc {
	return g.El("html", children...)
}

// Head returns an element with name "head" and the given children.
func Head(children ...g.Node) g.NodeFunc {
	return g.El("head", children...)
}

// Body returns an element with name "body" and the given children.
func Body(children ...g.Node) g.NodeFunc {
	return g.El("body", children...)
}

// Title returns an element with name "title" and a single Text child.
func Title(title string) g.NodeFunc {
	return g.El("title", g.Text(title))
}

func Meta(children ...g.Node) g.NodeFunc {
	return g.El("meta", children...)
}

func Link(children ...g.Node) g.NodeFunc {
	return g.El("link", children...)
}

func Style(children ...g.Node) g.NodeFunc {
	return g.El("style", children...)
}

func Base(children ...g.Node) g.NodeFunc {
	return g.El("base", children...)
}
