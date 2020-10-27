package el

import (
	g "github.com/maragudk/gomponents"
)

func Del(text string, children ...g.Node) g.NodeFunc {
	return g.El("del", g.Text(text), g.Group(children))
}

func Ins(text string, children ...g.Node) g.NodeFunc {
	return g.El("ins", g.Text(text), g.Group(children))
}
