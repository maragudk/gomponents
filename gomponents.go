// Package gomponents provides declarative view components in Go, that can render to HTML5.
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

// voidElements don't have end tags and must be treated differently in the rendering.
// See https://dev.w3.org/html5/spec-LC/syntax.html#void-elements
var voidElements = []string{"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr"}

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

// nodeType is for DOM Nodes that are either an element or an attribute.
type nodeType int

const (
	elementType = nodeType(iota)
	attributeType
)

// El creates an element DOM Node with a name and child Nodes.
// Use this if no convenience creator exists.
// See https://dev.w3.org/html5/spec-LC/syntax.html#elements-0 for how elements are rendered.
// No tags are ever omitted from normal tags, even though it's allowed for elements given at
// https://dev.w3.org/html5/spec-LC/syntax.html#optional-tags
// If an element is a void kind, non-attribute nodes are ignored.
func El(name string, children ...Node) NodeFunc {
	return func(w2 io.Writer) error {
		w := &statefulWriter{w: w2}

		w.Write([]byte("<" + name))

		for _, c := range children {
			renderChild(w, c, attributeType)
		}

		if isVoidKind(name) {
			w.Write([]byte(">"))
			return w.err
		}

		w.Write([]byte(">"))

		for _, c := range children {
			renderChild(w, c, elementType)
		}

		w.Write([]byte("</" + name + ">"))
		return w.err
	}
}

func isVoidKind(name string) bool {
	for _, e := range voidElements {
		if name == e {
			return true
		}
	}
	return false
}

// renderChild c to the given writer w if the node type is t.
func renderChild(w *statefulWriter, c Node, t nodeType) {
	if w.err != nil || c == nil {
		return
	}

	if g, ok := c.(group); ok {
		for _, groupC := range g.children {
			renderChild(w, groupC, t)
		}
		return
	}

	switch t {
	case elementType:
		if p, ok := c.(Placer); !ok || p.Place() == Outside {
			w.err = c.Render(w.w)
		}
	case attributeType:
		if p, ok := c.(Placer); ok && p.Place() == Inside {
			w.err = c.Render(w.w)
		}
	}
}

// statefulWriter only writes if no errors have occured earlier in its lifetime.
type statefulWriter struct {
	w   io.Writer
	err error
}

func (w *statefulWriter) Write(p []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(p)
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

func (g group) String() string {
	panic("cannot render group directly")
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

// Map something enumerable to a list of Nodes.
// Example:
// 	items := []string{"hat", "partyhat"}
//
// 	lis := g.Map(len(items), func(i int) g.Node {
// 		return g.El("li", g.Text(items[i]))
// 	})
//
// 	list := g.El("ul", lis...)
func Map(length int, cb func(i int) Node) []Node {
	var nodes []Node
	for i := 0; i < length; i++ {
		nodes = append(nodes, cb(i))
	}
	return nodes
}
