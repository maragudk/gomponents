package el

import (
	g "github.com/maragudk/gomponents"
)

func Br(children ...g.Node) g.NodeFunc {
	return g.El("br", children...)
}

func Hr(children ...g.Node) g.NodeFunc {
	return g.El("hr", children...)
}
