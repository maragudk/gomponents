package el

import (
	g "github.com/maragudk/gomponents"
)

func Embed(children ...g.Node) g.NodeFunc {
	return g.El("embed", children...)
}

func IFrame(children ...g.Node) g.NodeFunc {
	return g.El("iframe", children...)
}

func Object(children ...g.Node) g.NodeFunc {
	return g.El("object", children...)
}

func Param(children ...g.Node) g.NodeFunc {
	return g.El("param", children...)
}

func Picture(children ...g.Node) g.NodeFunc {
	return g.El("picture", children...)
}

func Source(children ...g.Node) g.NodeFunc {
	return g.El("source", children...)
}
