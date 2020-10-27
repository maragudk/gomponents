package el

import (
	g "github.com/maragudk/gomponents"
)

func Canvas(children ...g.Node) g.NodeFunc {
	return g.El("canvas", children...)
}

func NoScript(children ...g.Node) g.NodeFunc {
	return g.El("noscript", children...)
}

func Script(children ...g.Node) g.NodeFunc {
	return g.El("script", children...)
}
