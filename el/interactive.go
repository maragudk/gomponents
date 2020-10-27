package el

import (
	g "github.com/maragudk/gomponents"
)

func Details(children ...g.Node) g.NodeFunc {
	return g.El("details", children...)
}

func Dialog(children ...g.Node) g.NodeFunc {
	return g.El("dialog", children...)
}

func Menu(children ...g.Node) g.NodeFunc {
	return g.El("menu", children...)
}

func Summary(children ...g.Node) g.NodeFunc {
	return g.El("summary", children...)
}
