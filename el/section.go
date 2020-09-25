package el

import (
	g "github.com/maragudk/gomponents"
)

func H1(text string, children ...g.Node) g.NodeFunc {
	return g.El("h1", prepend(g.Text(text), children)...)
}

func H2(text string, children ...g.Node) g.NodeFunc {
	return g.El("h2", prepend(g.Text(text), children)...)
}

func H3(text string, children ...g.Node) g.NodeFunc {
	return g.El("h3", prepend(g.Text(text), children)...)
}

func H4(text string, children ...g.Node) g.NodeFunc {
	return g.El("h4", prepend(g.Text(text), children)...)
}

func H5(text string, children ...g.Node) g.NodeFunc {
	return g.El("h5", prepend(g.Text(text), children)...)
}

func H6(text string, children ...g.Node) g.NodeFunc {
	return g.El("h6", prepend(g.Text(text), children)...)
}
