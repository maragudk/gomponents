// Package el provides shortcuts and helpers to common HTML elements.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element for a list of elements.
package el

import (
	"strings"

	g "github.com/maragudk/gomponents"
)

// Document returns an special kind of Node that prefixes its children with the string "<!doctype html>".
func Document(children ...g.Node) g.NodeFunc {
	return func() string {
		var b strings.Builder
		b.WriteString("<!doctype html>")
		for _, c := range children {
			b.WriteString(c.Render())
		}
		return b.String()
	}
}

// HTML returns an element with name "html" and the given children.
func HTML(children ...g.Node) g.NodeFunc {
	return g.El("html", children...)
}

// Head returns an element with name "head" and the given children.
func Head(children ...g.Node) g.NodeFunc {
	return g.El("head", children...)
}

// Body returns an element with name "body" and the given children.
func Body(children ...g.Node) g.NodeFunc {
	return g.El("body", children...)
}

// Title returns an element with name "title" and a single Text child.
func Title(title string) g.NodeFunc {
	return g.El("title", g.Text(title))
}

func Meta(children ...g.Node) g.NodeFunc {
	return g.El("meta", children...)
}

func Link(children ...g.Node) g.NodeFunc {
	return g.El("link", children...)
}

func Style(children ...g.Node) g.NodeFunc {
	return g.El("style", children...)
}

func Div(children ...g.Node) g.NodeFunc {
	return g.El("div", children...)
}

func Span(children ...g.Node) g.NodeFunc {
	return g.El("span", children...)
}

func A(href string, children ...g.Node) g.NodeFunc {
	return g.El("a", prepend(g.Attr("href", href), children)...)
}

func P(children ...g.Node) g.NodeFunc {
	return g.El("p", children...)
}

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

func Ol(children ...g.Node) g.NodeFunc {
	return g.El("ol", children...)
}

func Ul(children ...g.Node) g.NodeFunc {
	return g.El("ul", children...)
}

func Li(children ...g.Node) g.NodeFunc {
	return g.El("li", children...)
}

func B(text string, children ...g.Node) g.NodeFunc {
	return g.El("b", prepend(g.Text(text), children)...)
}

func Strong(text string, children ...g.Node) g.NodeFunc {
	return g.El("strong", prepend(g.Text(text), children)...)
}

func I(text string, children ...g.Node) g.NodeFunc {
	return g.El("i", prepend(g.Text(text), children)...)
}

func Em(text string, children ...g.Node) g.NodeFunc {
	return g.El("em", prepend(g.Text(text), children)...)
}

func Img(src, alt string, children ...g.Node) g.NodeFunc {
	return g.El("img", prepend2(g.Attr("src", src), g.Attr("alt", alt), children)...)
}

func prepend(node g.Node, nodes []g.Node) []g.Node {
	newNodes := []g.Node{node}
	newNodes = append(newNodes, nodes...)
	return newNodes
}

func prepend2(node1, node2 g.Node, nodes []g.Node) []g.Node {
	newNodes := []g.Node{node1, node2}
	newNodes = append(newNodes, nodes...)
	return newNodes
}
