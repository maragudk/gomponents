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

// Classes returns an attribute with name "class" and the value being a sorted, space-separated string of all the keys,
// for which the corresponding value is true.
func Classes(cs map[string]bool) g.Node {
	var included []string
	for c, include := range cs {
		if include {
			included = append(included, c)
		}
	}
	sort.Strings(included)
	return g.Attr("class", strings.Join(included, " "))
}
