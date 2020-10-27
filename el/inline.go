package el

import (
	g "github.com/maragudk/gomponents"
)

func A(href string, children ...g.Node) g.NodeFunc {
	return g.El("a", g.Attr("href", href), g.Group(children))
}

func Abbr(text string, children ...g.Node) g.NodeFunc {
	return g.El("abbr", g.Text(text), g.Group(children))
}

func B(text string, children ...g.Node) g.NodeFunc {
	return g.El("b", g.Text(text), g.Group(children))
}

func Br(children ...g.Node) g.NodeFunc {
	return g.El("br", children...)
}

func Cite(children ...g.Node) g.NodeFunc {
	return g.El("cite", children...)
}

func Code(children ...g.Node) g.NodeFunc {
	return g.El("code", children...)
}

func Data(children ...g.Node) g.NodeFunc {
	return g.El("data", children...)
}

func Dfn(text string, children ...g.Node) g.NodeFunc {
	return g.El("dfn", g.Text(text), g.Group(children))
}

func Em(text string, children ...g.Node) g.NodeFunc {
	return g.El("em", g.Text(text), g.Group(children))
}

func I(text string, children ...g.Node) g.NodeFunc {
	return g.El("i", g.Text(text), g.Group(children))
}

func Kbd(text string, children ...g.Node) g.NodeFunc {
	return g.El("kbd", g.Text(text), g.Group(children))
}

func Mark(text string, children ...g.Node) g.NodeFunc {
	return g.El("mark", g.Text(text), g.Group(children))
}

func Q(text string, children ...g.Node) g.NodeFunc {
	return g.El("q", g.Text(text), g.Group(children))
}

func S(text string, children ...g.Node) g.NodeFunc {
	return g.El("s", g.Text(text), g.Group(children))
}

func Samp(text string, children ...g.Node) g.NodeFunc {
	return g.El("samp", g.Text(text), g.Group(children))
}

func Small(text string, children ...g.Node) g.NodeFunc {
	return g.El("small", g.Text(text), g.Group(children))
}

func Span(children ...g.Node) g.NodeFunc {
	return g.El("span", children...)
}

func Strong(text string, children ...g.Node) g.NodeFunc {
	return g.El("strong", g.Text(text), g.Group(children))
}

func Sub(text string, children ...g.Node) g.NodeFunc {
	return g.El("sub", g.Text(text), g.Group(children))
}

func Sup(text string, children ...g.Node) g.NodeFunc {
	return g.El("sup", g.Text(text), g.Group(children))
}

func Time(text string, children ...g.Node) g.NodeFunc {
	return g.El("time", g.Text(text), g.Group(children))
}

func U(text string, children ...g.Node) g.NodeFunc {
	return g.El("u", g.Text(text), g.Group(children))
}

func Var(text string, children ...g.Node) g.NodeFunc {
	return g.El("var", g.Text(text), g.Group(children))
}

func Wbr(children ...g.Node) g.NodeFunc {
	return g.El("wbr", children...)
}
