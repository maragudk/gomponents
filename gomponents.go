// Package gomponents provides components of DOM nodes for Go, that can render to an HTML Document.
// The primary interface is a Node, which has a single function Render, which should render
// the Node to a string. Furthermore, NodeFunc is a function which implements the Node interface
// by calling itself on Render.
// All DOM elements and attributes can be created by using the El and Attr functions.
// The package also provides a lot of convenience functions for creating elements and attributes
// with the most commonly used parameters. If they don't suffice, a fallback to El and Attr is always possible.
package gomponents

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// Node is a DOM node that can Render itself to a string representation.
type Node interface {
	Render() string
}

// NodeFunc is render function that is also a Node.
type NodeFunc func() string

func (n NodeFunc) Render() string {
	return n()
}

// String satisfies fmt.Stringer.
func (n NodeFunc) String() string {
	return n.Render()
}

// El creates an element DOM Node with a name and child Nodes.
// Use this if no convenience creator exists.
func El(name string, children ...Node) NodeFunc {
	return func() string {
		var b, attrString, childrenString strings.Builder

		b.WriteString("<")
		b.WriteString(name)

		if len(children) == 0 {
			b.WriteString("/>")
			return b.String()
		}

		for _, c := range children {
			if _, ok := c.(attr); ok {
				attrString.WriteString(c.Render())
			} else {
				childrenString.WriteString(c.Render())
			}
		}

		b.WriteString(attrString.String())

		if childrenString.Len() == 0 {
			b.WriteString("/>")
			return b.String()
		}

		b.WriteString(">")
		b.WriteString(childrenString.String())
		b.WriteString("</")
		b.WriteString(name)
		b.WriteString(">")
		return b.String()
	}
}

// Attr creates an attr DOM Node.
// If one parameter is passed, it's a name-only attribute (like "required").
// If two parameters are passed, it's a name-value attribute (like `class="header"`).
// More parameter counts make Attr panic.
// Use this if no convenience creator exists.
func Attr(name string, value ...string) Node {
	switch len(value) {
	case 0:
		return attr{name: name}
	case 1:
		return attr{name: name, value: &value[0]}
	default:
		panic("attribute must be just name or name and value pair")
	}
}

type attr struct {
	name  string
	value *string
}

func (a attr) Render() string {
	if a.value == nil {
		return fmt.Sprintf(" %v", a.name)
	}
	return fmt.Sprintf(` %v="%v"`, a.name, *a.value)
}

// String satisfies fmt.Stringer.
func (a attr) String() string {
	return a.Render()
}

// Text creates a text DOM Node that Renders the escaped string t.
func Text(t string) NodeFunc {
	return func() string {
		return template.HTMLEscaper(t)
	}
}

// Raw creates a raw Node that just Renders the unescaped string t.
func Raw(t string) NodeFunc {
	return func() string {
		return t
	}
}

// Write to the given io.Writer, returning any error.
func Write(w io.Writer, n Node) error {
	_, err := w.Write([]byte(n.Render()))
	return err
}
