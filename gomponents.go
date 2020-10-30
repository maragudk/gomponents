// Package gomponents provides declarative view components in Go, that can render to HTML.
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

// Node is a DOM node that can Render itself to a io.Writer.
type Node interface {
	Render(w io.Writer) error
}

// Placer can be implemented to tell Render functions where to place the string representation of a Node
// in the parent element.
type Placer interface {
	Place() Placement
}

// Placement is used with the Placer interface.
type Placement int

const (
	Outside = Placement(iota)
	Inside
)

// NodeFunc is render function that is also a Node.
type NodeFunc func(io.Writer) error

func (n NodeFunc) Render(w io.Writer) error {
	return n(w)
}

func (n NodeFunc) Place() Placement {
	return Outside
}

// String satisfies fmt.Stringer.
func (n NodeFunc) String() string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}

// El creates an element DOM Node with a name and child Nodes.
// Use this if no convenience creator exists.
func El(name string, children ...Node) NodeFunc {
	return func(w io.Writer) error {
		if _, err := w.Write([]byte("<" + name)); err != nil {
			return err
		}

		if len(children) == 0 {
			_, err := w.Write([]byte(" />"))
			return err
		}

		hasNonAttributeChild := false
		for _, c := range children {
			hasNonAttributeChild = hasNonAttributeChild || isNonAttributeChild(c)
			if err := renderAttr(w, c); err != nil {
				return err
			}
		}

		if !hasNonAttributeChild {
			_, err := w.Write([]byte(" />"))
			return err
		}

		if _, err := w.Write([]byte(">")); err != nil {
			return err
		}

		for _, c := range children {
			if err := renderChild(w, c); err != nil {
				return err
			}
		}

		_, err := w.Write([]byte("</" + name + ">"))
		return err
	}
}

func isNonAttributeChild(c Node) bool {
	if c == nil {
		return false
	}
	if g, ok := c.(group); ok {
		for _, groupC := range g.children {
			if isNonAttributeChild(groupC) {
				return true
			}
		}
		return false
	}
	if p, ok := c.(Placer); !ok || (ok && p.Place() == Outside) {
		return true
	}
	return false
}

func renderAttr(w io.Writer, c Node) error {
	if c == nil {
		return nil
	}
	if g, ok := c.(group); ok {
		for _, groupC := range g.children {
			if err := renderAttr(w, groupC); err != nil {
				return err
			}
		}
		return nil
	}
	if p, ok := c.(Placer); ok && p.Place() == Inside {
		return c.Render(w)
	}
	return nil
}

func renderChild(w io.Writer, c Node) error {
	if c == nil {
		return nil
	}
	if g, ok := c.(group); ok {
		for _, groupC := range g.children {
			if err := renderChild(w, groupC); err != nil {
				return err
			}
		}
		return nil
	}
	if p, ok := c.(Placer); ok && p.Place() == Inside {
		return nil
	}
	// If c doesn't implement Placer, default to outside
	return c.Render(w)
}

// Attr creates an attr DOM Node.
// If one parameter is passed, it's a name-only attribute (like "required").
// If two parameters are passed, it's a name-value attribute (like `class="header"`).
// More parameter counts make Attr panic.
// Use this if no convenience creator exists.
func Attr(name string, value ...string) Node {
	switch len(value) {
	case 0:
		return &attr{name: name}
	case 1:
		return &attr{name: name, value: &value[0]}
	default:
		panic("attribute must be just name or name and value pair")
	}
}

type attr struct {
	name  string
	value *string
}

func (a *attr) Render(w io.Writer) error {
	if a.value == nil {
		_, err := w.Write([]byte(" " + a.name))
		return err
	}
	_, err := w.Write([]byte(" " + a.name + `="` + *a.value + `"`))
	return err
}

func (a *attr) Place() Placement {
	return Inside
}

// String satisfies fmt.Stringer.
func (a *attr) String() string {
	var b strings.Builder
	_ = a.Render(&b)
	return b.String()
}

// Text creates a text DOM Node that Renders the escaped string t.
func Text(t string) NodeFunc {
	return func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(t)))
		return err
	}
}

// Textf creates a text DOM Node that Renders the interpolated and escaped string t.
func Textf(format string, a ...interface{}) NodeFunc {
	return func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(fmt.Sprintf(format, a...))))
		return err
	}
}

// Raw creates a raw Node that just Renders the unescaped string t.
func Raw(t string) NodeFunc {
	return func(w io.Writer) error {
		_, err := w.Write([]byte(t))
		return err
	}
}

type group struct {
	children []Node
}

func (g group) Render(io.Writer) error {
	panic("cannot render group directly")
}

// Group multiple Nodes into one Node. Useful for concatenation of Nodes in variadic functions.
// The resulting Node cannot Render directly, trying it will panic.
// Render must happen through a parent element created with El or a helper.
func Group(children []Node) Node {
	return group{children: children}
}
