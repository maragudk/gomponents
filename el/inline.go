package el

import (
	g "github.com/maragudk/gomponents"
)

func Span(children ...g.Node) g.NodeFunc {
	return g.El("span", children...)
}

func A(href string, children ...g.Node) g.NodeFunc {
	return g.El("a", prepend(g.Attr("href", href), children)...)
}

func B(text string, children ...g.Node) g.NodeFunc {
	return g.El("b", prepend(g.Text(text), children)...)
}

func Strong(text string, children ...g.Node) g.NodeFunc {
	return g.El("strong", prepend(g.Text(text), children)...)
}

func I(text string, children ...g.Node) g.NodeFunc {
	return g.El("i", prepend(g.Text(text), children)...)
}

func Em(text string, children ...g.Node) g.NodeFunc {
	return g.El("em", prepend(g.Text(text), children)...)
}
