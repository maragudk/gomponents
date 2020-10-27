package el

import (
	g "github.com/maragudk/gomponents"
)

func Area(children ...g.Node) g.NodeFunc {
	return g.El("area", children...)
}

func Audio(children ...g.Node) g.NodeFunc {
	return g.El("audio", children...)
}

func Img(src, alt string, children ...g.Node) g.NodeFunc {
	return g.El("img", g.Attr("src", src), g.Attr("alt", alt), g.Group(children))
}

func Map(children ...g.Node) g.NodeFunc {
	return g.El("map", children...)
}

func Track(children ...g.Node) g.NodeFunc {
	return g.El("track", children...)
}

func Video(children ...g.Node) g.NodeFunc {
	return g.El("video", children...)
}
