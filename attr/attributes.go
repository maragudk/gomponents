// Package attr provides shortcuts and helpers to common HTML attributes.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes for a list of attributes.
package attr

import (
	"sort"
	"strings"

	g "github.com/maragudk/gomponents"
)

// ID returns an attribute with name "id" and the given value.
func ID(v string) g.Node {
	return g.Attr("id", v)
}

// Class returns an attribute with name "class" and the given value.
func Class(v string) g.Node {
	return g.Attr("class", v)
}

// Classes is a map of strings to booleans, which Renders to an attribute with name "class".
// The attribute value is a sorted, space-separated string of all the map keys,
// for which the corresponding map value is true.
type Classes map[string]bool

func (c Classes) Render() string {
	var included []string
	for c, include := range c {
		if include {
			included = append(included, c)
		}
	}
	sort.Strings(included)
	return g.Attr("class", strings.Join(included, " ")).Render()
}

func (c Classes) Place() g.Placement {
	return g.Inside
}

// String satisfies fmt.Stringer.
func (c Classes) String() string {
	return c.Render()
}
