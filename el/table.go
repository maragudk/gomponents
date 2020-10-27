package el

import (
	g "github.com/maragudk/gomponents"
)

func Caption(text string, children ...g.Node) g.NodeFunc {
	return g.El("caption", g.Text(text), g.Group(children))
}

func Col(children ...g.Node) g.NodeFunc {
	return g.El("col", children...)
}

func ColGroup(children ...g.Node) g.NodeFunc {
	return g.El("colgroup", children...)
}

func Table(children ...g.Node) g.NodeFunc {
	return g.El("table", children...)
}

func TBody(children ...g.Node) g.NodeFunc {
	return g.El("tbody", children...)
}

func Td(children ...g.Node) g.NodeFunc {
	return g.El("td", children...)
}

func TFoot(children ...g.Node) g.NodeFunc {
	return g.El("tfoot", children...)
}

func Th(children ...g.Node) g.NodeFunc {
	return g.El("th", children...)
}

func THead(children ...g.Node) g.NodeFunc {
	return g.El("thead", children...)
}

func Tr(children ...g.Node) g.NodeFunc {
	return g.El("tr", children...)
}
