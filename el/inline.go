package el

import (
	g "github.com/maragudk/gomponents"
)

func Span(children ...g.Node) g.NodeFunc {
	return g.El("span", children...)
}

func A(href string, children ...g.Node) g.NodeFunc {
	return g.El("a", g.Attr("href", href), g.Group(children))
}

func B(text string, children ...g.Node) g.NodeFunc {
	return g.El("b", g.Text(text), g.Group(children))
}

func Strong(text string, children ...g.Node) g.NodeFunc {
	return g.El("strong", g.Text(text), g.Group(children))
}

func I(text string, children ...g.Node) g.NodeFunc {
	return g.El("i", g.Text(text), g.Group(children))
}

func Em(text string, children ...g.Node) g.NodeFunc {
	return g.El("em", g.Text(text), g.Group(children))
}
