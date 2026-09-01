// Package components provides high-level components and helpers that are composed of low-level elements and attributes.
package components

import (
	"bytes"
	"html"
	"io"
	"sort"
	"strings"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HTML5Props for [HTML5].
// Title is set no matter what, Description and Language elements only if the strings are non-empty.
type HTML5Props struct {
	Title       string
	Description string
	Language    string
	Head        g.Group
	Body        g.Group
	HTMLAttrs   g.Group
}

// HTML5 document template.
func HTML5(p HTML5Props) g.Node {
	return Doctype(
		HTML(g.If(p.Language != "", Lang(p.Language)), p.HTMLAttrs,
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				TitleEl(g.Text(p.Title)),
				g.If(p.Description != "", Meta(Name("description"), Content(p.Description))),
				p.Head,
			),
			Body(p.Body),
		),
	)
}

// Classes is a map of strings to booleans, which Renders to an attribute with name "class".
// The attribute value is a sorted, space-separated string of all the map keys,
// for which the corresponding map value is true.
type Classes map[string]bool

// Render satisfies [g.Node].
func (c Classes) Render(w io.Writer) error {
	included := make([]string, 0, len(c))
	for c, include := range c {
		if include {
			included = append(included, c)
		}
	}
	sort.Strings(included)
	return Class(strings.Join(included, " ")).Render(w)
}

func (c Classes) Type() g.NodeType {
	return g.AttributeType
}

// String satisfies [fmt.Stringer].
func (c Classes) String() string {
	var b strings.Builder
	_ = c.Render(&b)
	return b.String()
}

// JoinAttrs joins attributes with the given name on the first level of the given nodes.
// A [g.Group] is transparent and doesn't count as a level, so attributes inside one are joined
// at any depth, matching how groups are rendered.
// Attributes on non-direct descendants, for example inside a child element, are ignored.
// Non-empty attribute values are joined by spaces into a single attribute.
// Empty, whitespace-only, and boolean (valueless) attributes are deduplicated and discarded
// if a non-empty value exists. If only boolean/empty attributes match, a single boolean
// attribute is emitted.
// When both boolean and valued attributes match, the valued form takes precedence.
// The name is rendered unescaped and must be a trusted value, never user-controlled data.
// Note that this renders all first-level attributes, at any group depth, to check whether they
// should be processed.
func JoinAttrs(name string, children ...g.Node) g.Node {
	// The two shapes an attribute called name renders as. extractAttrValue compares against
	// them for every child, so build them once instead of per child.
	boolAttr := []byte(" " + name)
	attrPrefix := []byte(" " + name + `="`)

	// One buffer for every child. bytes.Buffer.Reset keeps its capacity, so inspecting the
	// second child doesn't re-grow it from nil the way a fresh strings.Builder would.
	var buf bytes.Buffer

	var attrValues []string
	result := make([]g.Node, 0, len(children))
	firstAttrIndex := -1
	sawBoolAttr := false

	// processNode checks a single child node and either collects its value or appends it to result.
	// Groups are unwrapped recursively, mirroring how they are rendered.
	var processNode func(n g.Node)
	processNode = func(n g.Node) {
		if group, ok := n.(g.Group); ok {
			for _, groupChild := range group {
				processNode(groupChild)
			}
			return
		}

		isGivenAttr, attrValue := extractAttrValue(&buf, boolAttr, attrPrefix, n)
		if !isGivenAttr {
			result = append(result, n)
			return
		}
		if attrValue == "" {
			sawBoolAttr = true
			if firstAttrIndex == -1 {
				firstAttrIndex = len(result)
				result = append(result, nil)
			}
			return
		}
		if attrValues == nil {
			// Only reached when there is something to join, so don't allocate for the
			// common case of a child list with no matching attribute at all.
			attrValues = make([]string, 0, len(children))
		}
		attrValues = append(attrValues, attrValue)
		if firstAttrIndex == -1 {
			firstAttrIndex = len(result)
			result = append(result, nil)
		}
	}

	for _, child := range children {
		processNode(child)
	}

	// If no matching attributes were found, just return the result now
	if firstAttrIndex == -1 {
		return g.Group(result)
	}

	// Insert joined attribute at the position of the first match
	if len(attrValues) > 0 {
		result[firstAttrIndex] = g.Attr(name, strings.Join(attrValues, " "))
	} else if sawBoolAttr {
		result[firstAttrIndex] = g.Attr(name)
	}
	return g.Group(result)
}

var quote = []byte(`"`)

type nodeTypeDescriber interface {
	Type() g.NodeType
}

func extractAttrValue(buf *bytes.Buffer, boolAttr, attrPrefix []byte, n g.Node) (bool, string) {
	// Ignore everything that is not an attribute
	if n, ok := n.(nodeTypeDescriber); !ok || n.Type() == g.ElementType {
		return false, ""
	}

	buf.Reset()
	if err := n.Render(buf); err != nil {
		return false, ""
	}

	// Only valid until the next Reset, so it must not outlive this call.
	rendered := buf.Bytes()

	// Match boolean attribute (e.g., ` required`)
	if bytes.Equal(rendered, boolAttr) {
		return true, ""
	}

	if !bytes.HasPrefix(rendered, attrPrefix) || !bytes.HasSuffix(rendered, quote) {
		return false, ""
	}

	// A node that rendered as just ` name="` has the prefix's own quote as its suffix and
	// nothing in between, so there is no value to take.
	if len(rendered) == len(attrPrefix) {
		return true, ""
	}

	// Unescape to get the original value, since it will be escaped again when the joined attribute is rendered
	v := html.UnescapeString(string(rendered[len(attrPrefix) : len(rendered)-1]))
	// Treat whitespace-only values the same as empty
	if strings.TrimSpace(v) == "" {
		return true, ""
	}
	return true, v
}
