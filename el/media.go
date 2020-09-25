package el

import (
	g "github.com/maragudk/gomponents"
)

func Img(src, alt string, children ...g.Node) g.NodeFunc {
	return g.El("img", prepend2(g.Attr("src", src), g.Attr("alt", alt), children)...)
}
