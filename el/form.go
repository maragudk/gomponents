package el

import (
	"fmt"

	g "github.com/maragudk/gomponents"
)

// Button returns an element with name "button" and the given children.
func Button(children ...g.Node) g.NodeFunc {
	return g.El("button", children...)
}

func Datalist(children ...g.Node) g.NodeFunc {
	return g.El("datalist", children...)
}

func Fieldset(children ...g.Node) g.NodeFunc {
	return g.El("fieldset", children...)
}

// Form returns an element with name "form", the given action and method attributes, and the given children.
func Form(action, method string, children ...g.Node) g.NodeFunc {
	return g.El("form", g.Attr("action", action), g.Attr("method", method), g.Group(children))
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

func Legend(children ...g.Node) g.NodeFunc {
	return g.El("legend", children...)
}

func Meter(children ...g.Node) g.NodeFunc {
	return g.El("meter", children...)
}

func OptGroup(children ...g.Node) g.NodeFunc {
	return g.El("optgroup", children...)
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
