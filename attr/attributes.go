// Package attr provides shortcuts and helpers to common HTML attributes.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes for a list of attributes.
package attr

import (
	"io"
	"sort"
	"strings"

	g "github.com/maragudk/gomponents"
)

// Classes is a map of strings to booleans, which Renders to an attribute with name "class".
// The attribute value is a sorted, space-separated string of all the map keys,
// for which the corresponding map value is true.
type Classes map[string]bool

func (c Classes) Render(w io.Writer) error {
	var included []string
	for c, include := range c {
		if include {
			included = append(included, c)
		}
	}
	sort.Strings(included)
	return g.Attr("class", strings.Join(included, " ")).Render(w)
}

func (c Classes) Place() g.Placement {
	return g.Inside
}

// String satisfies fmt.Stringer.
func (c Classes) String() string {
	var b strings.Builder
	_ = c.Render(&b)
	return b.String()
}
