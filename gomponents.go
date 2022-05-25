// Package gomponents provides view components in Go, that render to HTML 5.
//
// The primary interface is a Node. It describes a function Render, which should render the Node
// to the given writer as a string.
//
// All DOM elements and attributes can be created by using the El and Attr functions.
// The functions Text, Textf, and Raw can be used to create text nodes.
// See also helper functions Group, Map, and If.
//
// For basic HTML elements and attributes, see the package html.
// For higher-level HTML components, see the package components.
// For SVG elements and attributes, see the package svg.
// For HTTP helpers, see the package http.
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

// NodeType describes what type of Node it is, currently either an element or an attribute.
// Nodes default to being ElementType.
type NodeType int

const (
	ElementType = NodeType(iota)
	AttributeType
)

// nodeTypeDescriber can be implemented by Nodes to let callers know whether the Node is an ElementType or an AttributeType.
// This is used for rendering.
type nodeTypeDescriber interface {
	Type() NodeType
}

// NodeFunc is render function that is also a Node of ElementType.
type NodeFunc func(io.Writer) error

func (n NodeFunc) Render(w io.Writer) error {
	return n(w)
}

func (n NodeFunc) Type() NodeType {
	return ElementType
}

// String satisfies fmt.Stringer.
func (n NodeFunc) String() string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}

// El creates an element DOM Node with a name and child Nodes.
// See https://dev.w3.org/html5/spec-LC/syntax.html#elements-0 for how elements are rendered.
// No tags are ever omitted from normal tags, even though it's allowed for elements given at
// https://dev.w3.org/html5/spec-LC/syntax.html#optional-tags
// If an element is a void kind, non-attribute children nodes are ignored.
// Use this if no convenience creator exists.
func El(name string, children ...Node) Node {
	return NodeFunc(func(w2 io.Writer) error {
		w := &statefulWriter{w: w2}

		w.Write([]byte("<" + name))

		for _, c := range children {
			renderChild(w, c, AttributeType)
		}

		w.Write([]byte(">"))

		if isVoidKind(name) {
			return w.err
		}

		for _, c := range children {
			renderChild(w, c, ElementType)
		}

		w.Write([]byte("</" + name + ">"))
		return w.err
	})
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
func renderChild(w *statefulWriter, c Node, t NodeType) {
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
	case ElementType:
		if p, ok := c.(nodeTypeDescriber); !ok || p.Type() == ElementType {
			w.err = c.Render(w.w)
		}
	case AttributeType:
		if p, ok := c.(nodeTypeDescriber); ok && p.Type() == AttributeType {
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

// Attr creates an attribute DOM Node with a name and optional value.
// If only a name is passed, it's a name-only (boolean) attribute (like "required").
// If a name and value are passed, it's a name-value attribute (like `class="header"`).
// More than one value make Attr panic.
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
	_, err := w.Write([]byte(" " + a.name + `="` + template.HTMLEscapeString(*a.value) + `"`))
	return err
}

func (a *attr) Type() NodeType {
	return AttributeType
}

// String satisfies fmt.Stringer.
func (a *attr) String() string {
	var b strings.Builder
	_ = a.Render(&b)
	return b.String()
}

// Text creates a text DOM Node that Renders the escaped string t.
func Text(t string) Node {
	return NodeFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(t)))
		return err
	})
}

// Textf creates a text DOM Node that Renders the interpolated and escaped string t.
func Textf(format string, a ...interface{}) Node {
	return NodeFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(template.HTMLEscapeString(fmt.Sprintf(format, a...))))
		return err
	})
}

// Raw creates a text DOM Node that just Renders the unescaped string t.
func Raw(t string) Node {
	return NodeFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(t))
		return err
	})
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
func Map(length int, cb func(i int) Node) []Node {
	var nodes []Node
	for i := 0; i < length; i++ {
		nodes = append(nodes, cb(i))
	}
	return nodes
}

// If condition is true, return the given Node. Otherwise, return nil.
// This helper function is good for inlining elements conditionally.
func If(condition bool, n Node) Node {
	if condition {
		return n
	}
	return nil
}

// IfFunc is like If, except it takes a callback function that returns a Node instead of a Node directly.
// Great when the condition is a nil check and the function uses the value that's not nil.
func IfFunc(condition bool, cb func() Node) Node {
	if condition {
		return cb()
	}
	return nil
}
