package el

import (
	g "github.com/maragudk/gomponents"
)

func Img(src, alt string, children ...g.Node) g.NodeFunc {
	return g.El("img", g.Attr("src", src), g.Attr("alt", alt), g.Wrap(children...))
}
