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

func H1(text string) g.NodeFunc {
	return g.El("h1", g.Text(text))
}

func H2(text string) g.NodeFunc {
	return g.El("h2", g.Text(text))
}

func H3(text string) g.NodeFunc {
	return g.El("h3", g.Text(text))
}

func H4(text string) g.NodeFunc {
	return g.El("h4", g.Text(text))
}

func H5(text string) g.NodeFunc {
	return g.El("h5", g.Text(text))
}

func H6(text string) g.NodeFunc {
	return g.El("h6", g.Text(text))
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

func B(text string) g.NodeFunc {
	return g.El("b", g.Text(text))
}

func Strong(text string) g.NodeFunc {
	return g.El("strong", g.Text(text))
}

func I(text string) g.NodeFunc {
	return g.El("i", g.Text(text))
}

func Em(text string) g.NodeFunc {
	return g.El("em", g.Text(text))
}

func Img(src, alt string) g.NodeFunc {
	return g.El("img", g.Attr("src", src), g.Attr("alt", alt))
}

func prepend(node g.Node, nodes []g.Node) []g.Node {
	newNodes := []g.Node{node}
	newNodes = append(newNodes, nodes...)
	return newNodes
}
