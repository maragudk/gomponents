// Package svg provides common SVG elements and attributes.
// See https://developer.mozilla.org/en-US/docs/Web/SVG/Element for an overview.
package svg

import (
	g "github.com/maragudk/gomponents"
)

func Path(children ...g.Node) g.Node {
	return g.El("path", children...)
}

func SVG(children ...g.Node) g.Node {
	return g.El("svg", g.Attr("xmlns", "http://www.w3.org/2000/svg"), g.Group(children))
}
