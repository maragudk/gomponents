// Package components provides high-level components and helpers that are composed of low-level elements and attributes.
package components

import (
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
	Head        []g.Node
	Body        []g.Node
	HTMLAttrs   []g.Node
}

// HTML5 document template.
func HTML5(p HTML5Props) g.Node {
	return Doctype(
		HTML(g.If(p.Language != "", Lang(p.Language)), g.Group(p.HTMLAttrs),
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				TitleEl(g.Text(p.Title)),
				g.If(p.Description != "", Meta(Name("description"), Content(p.Description))),
				g.Group(p.Head),
			),
			Body(g.Group(p.Body)),
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

// JoinAttrs with the given name only on the first level of the given nodes.
// This means that attributes on non-direct descendants are ignored.
// Attribute values are joined by spaces.
// Note that this renders all first-level attributes to check whether they should be processed.
func JoinAttrs(name string, children ...g.Node) g.Node {
	var attrValues []string
	var result []g.Node
	firstAttrIndex := -1

	// Process all children
	for _, child := range children {
		// Handle groups explicitly because they may contain attributes
		if group, ok := child.(g.Group); ok {
			for _, groupChild := range group {
				isGivenAttr, attrValue := extractAttrValue(name, groupChild)
				if !isGivenAttr || attrValue == "" {
					result = append(result, groupChild)
					continue
				}

				attrValues = append(attrValues, attrValue)
				if firstAttrIndex == -1 {
					firstAttrIndex = len(result)
					result = append(result, nil)
				}
			}

			continue
		}

		// Handle non-group nodes essentially the same way
		isGivenAttr, attrValue := extractAttrValue(name, child)
		if !isGivenAttr || attrValue == "" {
			result = append(result, child)
			continue
		}

		attrValues = append(attrValues, attrValue)
		if firstAttrIndex == -1 {
			firstAttrIndex = len(result)
			result = append(result, nil)
		}
	}

	// If no attributes were found, just return the result now
	if firstAttrIndex == -1 {
		return g.Group(result)
	}

	// Insert joined attribute at the position of the first attribute
	result[firstAttrIndex] = g.Attr(name, strings.Join(attrValues, " "))
	return g.Group(result)
}

type nodeTypeDescriber interface {
	Type() g.NodeType
}

func extractAttrValue(name string, n g.Node) (bool, string) {
	// Ignore everything that is not an attribute
	if n, ok := n.(nodeTypeDescriber); !ok || n.Type() == g.ElementType {
		return false, ""
	}

	var b strings.Builder
	if err := n.Render(&b); err != nil {
		return false, ""
	}

	rendered := b.String()
	if !strings.HasPrefix(rendered, " "+name+`="`) || !strings.HasSuffix(rendered, `"`) {
		return false, ""
	}

	v := strings.TrimPrefix(rendered, " "+name+`="`)
	v = strings.TrimSuffix(v, `"`)
	return true, v
}
