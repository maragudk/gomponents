package el

import (
	g "github.com/maragudk/gomponents"
)

func BlockQuote(children ...g.Node) g.NodeFunc {
	return g.El("blockquote", children...)
}

func Dd(text string, children ...g.Node) g.NodeFunc {
	return g.El("dd", g.Text(text), g.Group(children))
}

func Div(children ...g.Node) g.NodeFunc {
	return g.El("div", children...)
}

func Dl(children ...g.Node) g.NodeFunc {
	return g.El("dl", children...)
}

func Dt(text string, children ...g.Node) g.NodeFunc {
	return g.El("dt", g.Text(text), g.Group(children))
}

func FigCaption(children ...g.Node) g.NodeFunc {
	return g.El("figcaption", children...)
}

func Figure(children ...g.Node) g.NodeFunc {
	return g.El("figure", children...)
}

func Hr(children ...g.Node) g.NodeFunc {
	return g.El("hr", children...)
}

func Li(children ...g.Node) g.NodeFunc {
	return g.El("li", children...)
}

func Ol(children ...g.Node) g.NodeFunc {
	return g.El("ol", children...)
}

func P(children ...g.Node) g.NodeFunc {
	return g.El("p", children...)
}

func Pre(children ...g.Node) g.NodeFunc {
	return g.El("pre", children...)
}

func Ul(children ...g.Node) g.NodeFunc {
	return g.El("ul", children...)
}
