// Package gomponents provides view components in Go, that render to HTML 5.
//
// The primary interface is a Node. It describes a function Render, which should render the Node
// to the given writer as a string.
//
// All DOM elements and attributes can be created by using the El and Attr functions.
// The functions Text, Textf, Raw, and Rawf can be used to create text nodes, either HTML-escaped or unescaped.
// See also helper functions Group, Map, and If for mapping data to Nodes and inserting them conditionally.
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

// Node is a DOM node that can Render itself to a io.Writer.
type Node interface {
	Render(w io.Writer) error
}

// NodeType describes what type of Node it is, currently either an ElementType or an AttributeType.
// This decides where a Node should be rendered.
// Nodes default to being ElementType.
type NodeType int

const (
	ElementType = NodeType(iota)
	AttributeType
)

// nodeTypeDescriber can be implemented by Nodes to let callers know whether the Node is
// an ElementType or an AttributeType. This is used for rendering.
type nodeTypeDescriber interface {
	Type() NodeType
}

// NodeFunc is a render function that is also a Node of ElementType.
type NodeFunc func(io.Writer) error

// Render satisfies Node.
func (n NodeFunc) Render(w io.Writer) error {
	return n(w)
}

// Type satisfies nodeTypeDescriber.
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
// If an element is a void element, non-attribute children nodes are ignored.
// Use this if no convenience creator exists.
func El(name string, children ...Node) Node {
	return NodeFunc(func(w2 io.Writer) error {
		w := &statefulWriter{w: w2}

		w.Write([]byte("<" + name))

		for _, c := range children {
			renderChild(w, c, AttributeType)
		}

		w.Write([]byte(">"))

		if isVoidElement(name) {
			return w.err
		}

		for _, c := range children {
			renderChild(w, c, ElementType)
		}

		w.Write([]byte("</" + name + ">"))
		return w.err
	})
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

// statefulWriter only writes if no errors have occurred earlier in its lifetime.
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

// voidElements don't have end tags and must be treated differently in the rendering.
// See https://dev.w3.org/html5/spec-LC/syntax.html#void-elements
var voidElements = map[string]struct{}{
	"area":    {},
	"base":    {},
	"br":      {},
	"col":     {},
	"command": {},
	"embed":   {},
	"hr":      {},
	"img":     {},
	"input":   {},
	"keygen":  {},
	"link":    {},
	"meta":    {},
	"param":   {},
	"source":  {},
	"track":   {},
	"wbr":     {},
}

func isVoidElement(name string) bool {
	_, ok := voidElements[name]
	return ok
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

// Render satisfies Node.
func (a *attr) Render(w io.Writer) error {
	if a.value == nil {
		_, err := w.Write([]byte(" " + a.name))
		return err
	}
	_, err := w.Write([]byte(" " + a.name + `="` + template.HTMLEscapeString(*a.value) + `"`))
	return err
}

// Type satisfies nodeTypeDescriber.
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

// Textf creates a text DOM Node that Renders the interpolated and escaped string format.
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

// Rawf creates a text DOM Node that just Renders the interpolated and unescaped string format.
func Rawf(format string, a ...interface{}) Node {
	return NodeFunc(func(w io.Writer) error {
		_, err := w.Write([]byte(fmt.Sprintf(format, a...)))
		return err
	})
}

type group struct {
	children []Node
}

// String satisfies fmt.Stringer.
func (g group) String() string {
	panic("cannot render group directly")
}

// Render satisfies Node.
func (g group) Render(io.Writer) error {
	panic("cannot render group directly")
}

// Group multiple Nodes into one Node. Useful for concatenation of Nodes in variadic functions.
// The resulting Node cannot Render directly, trying it will panic.
// Render must happen through a parent element created with El or a helper.
func Group(children []Node) Node {
	return group{children: children}
}

// If condition is true, return the given Node. Otherwise, return nil.
// This helper function is good for inlining elements conditionally.
// If your condition and node involve a nilable variable, use iff because
// go will evaluate the node regardless of the condition.
func If(condition bool, n Node) Node {
	if condition {
		return n
	}
	return nil
}

// Iff execute the function f if condition is true, otherwise return nil.
// it is the preferred way to conditionally render a node if the node involves a nilable variable.
func Iff(condition bool, f func() Node) Node {
	if condition {
		return f()
	}
	return nil
}
