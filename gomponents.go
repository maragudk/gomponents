// Package gomponents provides HTML components in Go, that render to HTML 5.
//
// The primary interface is a [Node]. It defines a function Render, which should render the [Node]
// to the given writer as a string.
//
// All DOM elements and attributes can be created by using the [El] and [Attr] functions.
//
// The functions [Text], [Textf], [Raw], and [Rawf] can be used to create text nodes, either HTML-escaped or unescaped.
//
// See also helper functions [Map], [If], and [Iff] for mapping data to nodes and inserting them conditionally.
//
// There's also the [Group] type, which is a slice of [Node]-s that can be rendered as one [Node].
//
// For basic HTML elements and attributes, see the package html.
//
// For higher-level HTML components, see the package components.
//
// For HTTP helpers, see the package http.
package gomponents

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// Node is a DOM node that can Render itself to a [io.Writer].
type Node interface {
	Render(w io.Writer) error
}

// NodeType describes what type of [Node] it is, currently either an [ElementType] or an [AttributeType].
// This decides where a [Node] should be rendered.
// Nodes default to being [ElementType].
type NodeType int

const (
	ElementType = NodeType(iota)
	AttributeType
)

// nodeTypeDescriber can be implemented by Nodes to let callers know whether the [Node] is
// an [ElementType] or an [AttributeType].
// See [NodeType].
type nodeTypeDescriber interface {
	Type() NodeType
}

// NodeFunc is a render function that is also a [Node] of [ElementType].
type NodeFunc func(io.Writer) error

// Render satisfies [Node].
func (n NodeFunc) Render(w io.Writer) error {
	return n(w)
}

// Type satisfies [nodeTypeDescriber].
func (n NodeFunc) Type() NodeType {
	return ElementType
}

// String satisfies [fmt.Stringer].
func (n NodeFunc) String() string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}

var (
	lt      = []byte("<")
	gt      = []byte(">")
	ltSlash = []byte("</")
)

// El creates an element DOM [Node] with a name and child Nodes.
// See https://dev.w3.org/html5/spec-LC/syntax.html#elements-0 for how elements are rendered.
// No tags are ever omitted from normal tags, even though it's allowed for elements given at
// https://dev.w3.org/html5/spec-LC/syntax.html#optional-tags
// If an element is a void element, non-attribute children nodes are ignored.
// Use this if no convenience creator exists in the html package.
func El(name string, children ...Node) Node {
	return NodeFunc(func(w io.Writer) error {
		var err error

		sw, ok := w.(io.StringWriter)

		if _, err = w.Write(lt); err != nil {
			return err
		}

		if ok {
			if _, err = sw.WriteString(name); err != nil {
				return err
			}
		} else {
			if _, err = w.Write([]byte(name)); err != nil {
				return err
			}
		}

		for _, c := range children {
			if err = renderChild(w, c, AttributeType); err != nil {
				return err
			}
		}

		if _, err = w.Write(gt); err != nil {
			return err
		}

		if isVoidElement(name) {
			return nil
		}

		for _, c := range children {
			if err = renderChild(w, c, ElementType); err != nil {
				return err
			}
		}

		if _, err = w.Write(ltSlash); err != nil {
			return err
		}

		if ok {
			if _, err = sw.WriteString(name); err != nil {
				return err
			}
		} else {
			if _, err = w.Write([]byte(name)); err != nil {
				return err
			}
		}

		if _, err = w.Write(gt); err != nil {
			return err
		}

		return nil
	})
}

// renderChild c to the given writer w if the node type is desiredType.
func renderChild(w io.Writer, c Node, desiredType NodeType) error {
	if c == nil {
		return nil
	}

	// Rendering groups like this is still important even though a group can render itself,
	// since otherwise attributes will sometimes be ignored.
	if g, ok := c.(Group); ok {
		for _, groupC := range g {
			if err := renderChild(w, groupC, desiredType); err != nil {
				return err
			}
		}
		return nil
	}

	switch desiredType {
	case ElementType:
		if p, ok := c.(nodeTypeDescriber); !ok || p.Type() == desiredType {
			if err := c.Render(w); err != nil {
				return err
			}
		}
	case AttributeType:
		if p, ok := c.(nodeTypeDescriber); ok && p.Type() == desiredType {
			if err := c.Render(w); err != nil {
				return err
			}
		}
	}

	return nil
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

var (
	space      = []byte(" ")
	equalQuote = []byte(`="`)
	quote      = []byte(`"`)
)

// Attr creates an attribute DOM [Node] with a name and optional value.
// If only a name is passed, it's a name-only (boolean) attribute (like "required").
// If a name and value are passed, it's a name-value attribute (like `class="header"`).
// More than one value make [Attr] panic.
// Use this if no convenience creator exists in the html package.
func Attr(name string, value ...string) Node {
	if len(value) > 1 {
		panic("attribute must be just name or name and value pair")
	}

	return attrFunc(func(w io.Writer) error {
		var err error

		sw, ok := w.(io.StringWriter)

		if _, err = w.Write(space); err != nil {
			return err
		}

		// Attribute name
		if ok {
			if _, err = sw.WriteString(name); err != nil {
				return err
			}
		} else {
			if _, err = w.Write([]byte(name)); err != nil {
				return err
			}
		}

		if len(value) == 0 {
			return nil
		}

		if _, err = w.Write(equalQuote); err != nil {
			return err
		}

		// Attribute value
		if ok {
			if _, err = sw.WriteString(template.HTMLEscapeString(value[0])); err != nil {
				return err
			}
		} else {
			if _, err = w.Write([]byte(template.HTMLEscapeString(value[0]))); err != nil {
				return err
			}
		}

		if _, err = w.Write(quote); err != nil {
			return err
		}

		return nil
	})
}

// attrFunc is a render function that is also a [Node] of [AttributeType].
// It's basically the same as [NodeFunc], but for attributes.
type attrFunc func(io.Writer) error

// Render satisfies [Node].
func (a attrFunc) Render(w io.Writer) error {
	return a(w)
}

// Type satisfies [nodeTypeDescriber].
func (a attrFunc) Type() NodeType {
	return AttributeType
}

// String satisfies [fmt.Stringer].
func (a attrFunc) String() string {
	var b strings.Builder
	_ = a.Render(&b)
	return b.String()
}

// Text creates a text DOM [Node] that Renders the escaped string t.
func Text(t string) Node {
	return NodeFunc(func(w io.Writer) error {
		if w, ok := w.(io.StringWriter); ok {
			_, err := w.WriteString(template.HTMLEscapeString(t))
			return err
		}
		_, err := w.Write([]byte(template.HTMLEscapeString(t)))
		return err
	})
}

// Textf creates a text DOM [Node] that Renders the interpolated and escaped string format.
func Textf(format string, a ...interface{}) Node {
	return NodeFunc(func(w io.Writer) error {
		if w, ok := w.(io.StringWriter); ok {
			_, err := w.WriteString(template.HTMLEscapeString(fmt.Sprintf(format, a...)))
			return err
		}
		_, err := w.Write([]byte(template.HTMLEscapeString(fmt.Sprintf(format, a...))))
		return err
	})
}

// Raw creates a text DOM [Node] that just Renders the unescaped string t.
func Raw(t string) Node {
	return NodeFunc(func(w io.Writer) error {
		if w, ok := w.(io.StringWriter); ok {
			_, err := w.WriteString(t)
			return err
		}
		_, err := w.Write([]byte(t))
		return err
	})
}

// Rawf creates a text DOM [Node] that just Renders the interpolated and unescaped string format.
func Rawf(format string, a ...interface{}) Node {
	return NodeFunc(func(w io.Writer) error {
		if w, ok := w.(io.StringWriter); ok {
			_, err := w.WriteString(fmt.Sprintf(format, a...))
			return err
		}
		_, err := fmt.Fprintf(w, format, a...)
		return err
	})
}

// Map a slice of anything to a [Group] (which is just a slice of [Node]-s).
func Map[T any](ts []T, cb func(T) Node) Group {
	nodes := make([]Node, 0, len(ts))
	for _, t := range ts {
		nodes = append(nodes, cb(t))
	}
	return nodes
}

// Group a slice of [Node]-s into one Node, while still being usable like a regular slice of [Node]-s.
// A [Group] can render directly, but if any of the direct children are [AttributeType], they will be ignored,
// to not produce invalid HTML.
type Group []Node

// String satisfies [fmt.Stringer].
func (g Group) String() string {
	var b strings.Builder
	_ = g.Render(&b)
	return b.String()
}

// Render satisfies [Node].
func (g Group) Render(w io.Writer) error {
	for _, c := range g {
		if err := renderChild(w, c, ElementType); err != nil {
			return err
		}
	}
	return nil
}

// If condition is true, return the given [Node]. Otherwise, return nil.
// This helper function is good for inlining elements conditionally.
// If it's important that the given [Node] is only evaluated if condition is true
// (for example, when using nilable variables), use [Iff] instead.
func If(condition bool, n Node) Node {
	if condition {
		return n
	}
	return nil
}

// Iff condition is true, call the given function. Otherwise, return nil.
// This helper function is good for inlining elements conditionally when the node depends on nilable data,
// or some other code that could potentially panic.
// If you just need simple conditional rendering, see [If].
func Iff(condition bool, f func() Node) Node {
	if condition {
		return f()
	}
	return nil
}
