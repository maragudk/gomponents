package el

import (
	g "github.com/maragudk/gomponents"
)

// Address returns an element with name "address" and the given children.
func Address(children ...g.Node) g.NodeFunc {
	return g.El("address", children...)
}

// Article returns an element with name "article" and the given children.
func Article(children ...g.Node) g.NodeFunc {
	return g.El("article", children...)
}

// Aside returns an element with name "aside" and the given children.
func Aside(children ...g.Node) g.NodeFunc {
	return g.El("aside", children...)
}

// Footer returns an element with name "footer" and the given children.
func Footer(children ...g.Node) g.NodeFunc {
	return g.El("footer", children...)
}

// Header returns an element with name "header" and the given children.
func Header(children ...g.Node) g.NodeFunc {
	return g.El("header", children...)
}

// H1 returns an element with name "h1", the given text content, and the given children.
func H1(text string, children ...g.Node) g.NodeFunc {
	return g.El("h1", g.Text(text), g.Wrap(children...))
}

// H2 returns an element with name "h2", the given text content, and the given children.
func H2(text string, children ...g.Node) g.NodeFunc {
	return g.El("h2", g.Text(text), g.Wrap(children...))
}

// H3 returns an element with name "h3", the given text content, and the given children.
func H3(text string, children ...g.Node) g.NodeFunc {
	return g.El("h3", g.Text(text), g.Wrap(children...))
}

// H4 returns an element with name "h4", the given text content, and the given children.
func H4(text string, children ...g.Node) g.NodeFunc {
	return g.El("h4", g.Text(text), g.Wrap(children...))
}

// H5 returns an element with name "h5", the given text content, and the given children.
func H5(text string, children ...g.Node) g.NodeFunc {
	return g.El("h5", g.Text(text), g.Wrap(children...))
}

// H6 returns an element with name "h6", the given text content, and the given children.
func H6(text string, children ...g.Node) g.NodeFunc {
	return g.El("h6", g.Text(text), g.Wrap(children...))
}

// HGroup returns an element with name "hgroup" and the given children.
func HGroup(children ...g.Node) g.NodeFunc {
	return g.El("hgroup", children...)
}

// Main returns an element with name "main" and the given children.
func Main(children ...g.Node) g.NodeFunc {
	return g.El("main", children...)
}

// Nav returns an element with name "nav" and the given children.
func Nav(children ...g.Node) g.NodeFunc {
	return g.El("nav", children...)
}

// Section returns an element with name "section" and the given children.
func Section(children ...g.Node) g.NodeFunc {
	return g.El("section", children...)
}
