package svg

import (
	g "github.com/maragudk/gomponents"
)

func ClipRule(v string) g.Node {
	return g.Attr("clip-rule", v)
}

func D(v string) g.Node {
	return g.Attr("d", v)
}

func Fill(v string) g.Node {
	return g.Attr("fill", v)
}

func FillRule(v string) g.Node {
	return g.Attr("fill-rule", v)
}

func Stroke(v string) g.Node {
	return g.Attr("stroke", v)
}

func StrokeWidth(v string) g.Node {
	return g.Attr("stroke-width", v)
}

func ViewBox(v string) g.Node {
	return g.Attr("viewBox", v)
}
