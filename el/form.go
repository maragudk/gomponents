package el

import (
	"fmt"

	g "github.com/maragudk/gomponents"
)

// Button returns an element with name "button" and the given children.
func Button(children ...g.Node) g.NodeFunc {
	return g.El("button", children...)
}

// Form returns an element with name "form", the given action and method attributes, and the given children.
func Form(action, method string, children ...g.Node) g.NodeFunc {
	return g.El("form", prepend2(g.Attr("action", action), g.Attr("method", method), children)...)
}

// Input returns an element with name "input", the given type and name attributes, and the given children.
// Note that "type" is a keyword in Go, so the parameter is called typ.
func Input(typ, name string, children ...g.Node) g.NodeFunc {
	return g.El("input", prepend2(g.Attr("type", typ), g.Attr("name", name), children)...)
}

// Label returns an element with name "label", the given for attribute, and the given children.
// Note that "for" is a keyword in Go, so the parameter is called forr.
func Label(forr string, children ...g.Node) g.NodeFunc {
	return g.El("label", prepend(g.Attr("for", forr), children)...)
}

// Option returns an element with name "option", the given text content and value attribute, and the given children.
func Option(text, value string, children ...g.Node) g.NodeFunc {
	return g.El("option", prepend2(g.Attr("value", value), g.Text(text), children)...)
}

// Progress returns an element with name "progress", the given value and max attributes, and the given children.
func Progress(value, max float64, children ...g.Node) g.NodeFunc {
	return g.El("progress", prepend2(
		g.Attr("value", fmt.Sprintf("%v", value)),
		g.Attr("max", fmt.Sprintf("%v", max)),
		children)...)
}

// Select returns an element with name "select", the given name attribute, and the given children.
func Select(name string, children ...g.Node) g.NodeFunc {
	return g.El("select", prepend(g.Attr("name", name), children)...)
}

// Textarea returns an element with name "textarea", the given name attribute, and the given children.
func Textarea(name string, children ...g.Node) g.NodeFunc {
	return g.El("textarea", prepend(g.Attr("name", name), children)...)
}
