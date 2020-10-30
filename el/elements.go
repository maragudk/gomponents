// Package el provides shortcuts and helpers to common HTML elements.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element for a list of elements.
package el

import (
	"fmt"
	"io"

	g "github.com/maragudk/gomponents"
)

func A(href string, children ...g.Node) g.NodeFunc {
	return g.El("a", g.Attr("href", href), g.Group(children))
}

// Document returns an special kind of Node that prefixes its child with the string "<!doctype html>".
func Document(child g.Node) g.NodeFunc {
	return func(w io.Writer) error {
		if _, err := w.Write([]byte("<!doctype html>")); err != nil {
			return err
		}
		return child.Render(w)
	}
}

// Form returns an element with name "form", the given action and method attributes, and the given children.
func Form(action, method string, children ...g.Node) g.NodeFunc {
	return g.El("form", g.Attr("action", action), g.Attr("method", method), g.Group(children))
}

func Img(src, alt string, children ...g.Node) g.NodeFunc {
	return g.El("img", g.Attr("src", src), g.Attr("alt", alt), g.Group(children))
}

// Input returns an element with name "input", the given type and name attributes, and the given children.
// Note that "type" is a keyword in Go, so the parameter is called typ.
func Input(typ, name string, children ...g.Node) g.NodeFunc {
	return g.El("input", g.Attr("type", typ), g.Attr("name", name), g.Group(children))
}

// Label returns an element with name "label", the given for attribute, and the given children.
// Note that "for" is a keyword in Go, so the parameter is called forr.
func Label(forr string, children ...g.Node) g.NodeFunc {
	return g.El("label", g.Attr("for", forr), g.Group(children))
}

// Option returns an element with name "option", the given text content and value attribute, and the given children.
func Option(text, value string, children ...g.Node) g.NodeFunc {
	return g.El("option", g.Attr("value", value), g.Text(text), g.Group(children))
}

// Progress returns an element with name "progress", the given value and max attributes, and the given children.
func Progress(value, max float64, children ...g.Node) g.NodeFunc {
	return g.El("progress",
		g.Attr("value", fmt.Sprintf("%v", value)),
		g.Attr("max", fmt.Sprintf("%v", max)),
		g.Group(children))
}

// Select returns an element with name "select", the given name attribute, and the given children.
func Select(name string, children ...g.Node) g.NodeFunc {
	return g.El("select", g.Attr("name", name), g.Group(children))
}

// Textarea returns an element with name "textarea", the given name attribute, and the given children.
func Textarea(name string, children ...g.Node) g.NodeFunc {
	return g.El("textarea", g.Attr("name", name), g.Group(children))
}
