package el

import (
	g "github.com/maragudk/gomponents"
)

func Abbr(text string, children ...g.Node) g.NodeFunc {
	return g.El("abbr", g.Text(text), g.Group(children))
}

func B(text string, children ...g.Node) g.NodeFunc {
	return g.El("b", g.Text(text), g.Group(children))
}

func Caption(text string, children ...g.Node) g.NodeFunc {
	return g.El("caption", g.Text(text), g.Group(children))
}

func Dd(text string, children ...g.Node) g.NodeFunc {
	return g.El("dd", g.Text(text), g.Group(children))
}

func Del(text string, children ...g.Node) g.NodeFunc {
	return g.El("del", g.Text(text), g.Group(children))
}

func Dfn(text string, children ...g.Node) g.NodeFunc {
	return g.El("dfn", g.Text(text), g.Group(children))
}

func Dt(text string, children ...g.Node) g.NodeFunc {
	return g.El("dt", g.Text(text), g.Group(children))
}

func Em(text string, children ...g.Node) g.NodeFunc {
	return g.El("em", g.Text(text), g.Group(children))
}

func FigCaption(text string, children ...g.Node) g.NodeFunc {
	return g.El("figcaption", g.Text(text), g.Group(children))
}

// H1 returns an element with name "h1", the given text content, and the given children.
func H1(text string, children ...g.Node) g.NodeFunc {
	return g.El("h1", g.Text(text), g.Group(children))
}

// H2 returns an element with name "h2", the given text content, and the given children.
func H2(text string, children ...g.Node) g.NodeFunc {
	return g.El("h2", g.Text(text), g.Group(children))
}

// H3 returns an element with name "h3", the given text content, and the given children.
func H3(text string, children ...g.Node) g.NodeFunc {
	return g.El("h3", g.Text(text), g.Group(children))
}

// H4 returns an element with name "h4", the given text content, and the given children.
func H4(text string, children ...g.Node) g.NodeFunc {
	return g.El("h4", g.Text(text), g.Group(children))
}

// H5 returns an element with name "h5", the given text content, and the given children.
func H5(text string, children ...g.Node) g.NodeFunc {
	return g.El("h5", g.Text(text), g.Group(children))
}

// H6 returns an element with name "h6", the given text content, and the given children.
func H6(text string, children ...g.Node) g.NodeFunc {
	return g.El("h6", g.Text(text), g.Group(children))
}

func I(text string, children ...g.Node) g.NodeFunc {
	return g.El("i", g.Text(text), g.Group(children))
}

func Ins(text string, children ...g.Node) g.NodeFunc {
	return g.El("ins", g.Text(text), g.Group(children))
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

func Title(title string, children ...g.Node) g.NodeFunc {
	return g.El("title", g.Text(title), g.Group(children))
}

func U(text string, children ...g.Node) g.NodeFunc {
	return g.El("u", g.Text(text), g.Group(children))
}

func Var(text string, children ...g.Node) g.NodeFunc {
	return g.El("var", g.Text(text), g.Group(children))
}
