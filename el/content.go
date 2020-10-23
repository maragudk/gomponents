package el

import (
	g "github.com/maragudk/gomponents"
)

func Div(children ...g.Node) g.NodeFunc {
	return g.El("div", children...)
}

func Ol(children ...g.Node) g.NodeFunc {
	return g.El("ol", children...)
}

func Ul(children ...g.Node) g.NodeFunc {
	return g.El("ul", children...)
}

func Li(children ...g.Node) g.NodeFunc {
	return g.El("li", children...)
}

func P(children ...g.Node) g.NodeFunc {
	return g.El("p", children...)
}

func Br(children ...g.Node) g.NodeFunc {
	return g.El("br", children...)
}

func Hr(children ...g.Node) g.NodeFunc {
	return g.El("hr", children...)
}
