package attr

import (
	g "github.com/maragudk/gomponents"
)

// Placeholder returns an attribute with name "placeholder" and the given value.
func Placeholder(v string) g.Node {
	return g.Attr("placeholder", v)
}

// Required returns an attribute with name "required".
func Required() g.Node {
	return g.Attr("required")
}
