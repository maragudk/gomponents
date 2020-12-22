// Package html provides common HTML elements and attributes.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element for a list of elements.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes for a list of attributes.
package html

import (
	"io"

	g "github.com/maragudk/gomponents"
)

// Doctype returns a special kind of Node that prefixes its sibling with the string "<!doctype html>".
func Doctype(sibling g.Node) g.NodeFunc {
	return func(w io.Writer) error {
		if _, err := w.Write([]byte("<!doctype html>")); err != nil {
			return err
		}
		return sibling.Render(w)
	}
}

func A(children ...g.Node) g.NodeFunc {
	return g.El("a", children...)
}

func Address(children ...g.Node) g.NodeFunc {
	return g.El("address", children...)
}

func Area(children ...g.Node) g.NodeFunc {
	return g.El("area", children...)
}

func Article(children ...g.Node) g.NodeFunc {
	return g.El("article", children...)
}

func Aside(children ...g.Node) g.NodeFunc {
	return g.El("aside", children...)
}

func Audio(children ...g.Node) g.NodeFunc {
	return g.El("audio", children...)
}

func Base(children ...g.Node) g.NodeFunc {
	return g.El("base", children...)
}

func BlockQuote(children ...g.Node) g.NodeFunc {
	return g.El("blockquote", children...)
}

func Body(children ...g.Node) g.NodeFunc {
	return g.El("body", children...)
}

func Br(children ...g.Node) g.NodeFunc {
	return g.El("br", children...)
}

func Button(children ...g.Node) g.NodeFunc {
	return g.El("button", children...)
}

func Canvas(children ...g.Node) g.NodeFunc {
	return g.El("canvas", children...)
}

func Cite(children ...g.Node) g.NodeFunc {
	return g.El("cite", children...)
}

func Code(children ...g.Node) g.NodeFunc {
	return g.El("code", children...)
}

func Col(children ...g.Node) g.NodeFunc {
	return g.El("col", children...)
}

func ColGroup(children ...g.Node) g.NodeFunc {
	return g.El("colgroup", children...)
}

func DataEl(children ...g.Node) g.NodeFunc {
	return g.El("data", children...)
}

func DataList(children ...g.Node) g.NodeFunc {
	return g.El("datalist", children...)
}

func Details(children ...g.Node) g.NodeFunc {
	return g.El("details", children...)
}

func Dialog(children ...g.Node) g.NodeFunc {
	return g.El("dialog", children...)
}

func Div(children ...g.Node) g.NodeFunc {
	return g.El("div", children...)
}

func Dl(children ...g.Node) g.NodeFunc {
	return g.El("dl", children...)
}

func Embed(children ...g.Node) g.NodeFunc {
	return g.El("embed", children...)
}

func FormEl(children ...g.Node) g.NodeFunc {
	return g.El("form", children...)
}

func FieldSet(children ...g.Node) g.NodeFunc {
	return g.El("fieldset", children...)
}

func Figure(children ...g.Node) g.NodeFunc {
	return g.El("figure", children...)
}

func Footer(children ...g.Node) g.NodeFunc {
	return g.El("footer", children...)
}

func Head(children ...g.Node) g.NodeFunc {
	return g.El("head", children...)
}

func Header(children ...g.Node) g.NodeFunc {
	return g.El("header", children...)
}

func HGroup(children ...g.Node) g.NodeFunc {
	return g.El("hgroup", children...)
}

func Hr(children ...g.Node) g.NodeFunc {
	return g.El("hr", children...)
}

func HTML(children ...g.Node) g.NodeFunc {
	return g.El("html", children...)
}

func IFrame(children ...g.Node) g.NodeFunc {
	return g.El("iframe", children...)
}

func Img(children ...g.Node) g.NodeFunc {
	return g.El("img", children...)
}

func Input(children ...g.Node) g.NodeFunc {
	return g.El("input", children...)
}

func Label(children ...g.Node) g.NodeFunc {
	return g.El("label", children...)
}

func Legend(children ...g.Node) g.NodeFunc {
	return g.El("legend", children...)
}

func Li(children ...g.Node) g.NodeFunc {
	return g.El("li", children...)
}

func Link(children ...g.Node) g.NodeFunc {
	return g.El("link", children...)
}

func Main(children ...g.Node) g.NodeFunc {
	return g.El("main", children...)
}

func Menu(children ...g.Node) g.NodeFunc {
	return g.El("menu", children...)
}

func Meta(children ...g.Node) g.NodeFunc {
	return g.El("meta", children...)
}

func Meter(children ...g.Node) g.NodeFunc {
	return g.El("meter", children...)
}

func Nav(children ...g.Node) g.NodeFunc {
	return g.El("nav", children...)
}

func NoScript(children ...g.Node) g.NodeFunc {
	return g.El("noscript", children...)
}

func Object(children ...g.Node) g.NodeFunc {
	return g.El("object", children...)
}

func Ol(children ...g.Node) g.NodeFunc {
	return g.El("ol", children...)
}

func OptGroup(children ...g.Node) g.NodeFunc {
	return g.El("optgroup", children...)
}

func Option(children ...g.Node) g.NodeFunc {
	return g.El("option", children...)
}

func P(children ...g.Node) g.NodeFunc {
	return g.El("p", children...)
}

func Param(children ...g.Node) g.NodeFunc {
	return g.El("param", children...)
}

func Picture(children ...g.Node) g.NodeFunc {
	return g.El("picture", children...)
}

func Pre(children ...g.Node) g.NodeFunc {
	return g.El("pre", children...)
}

func Progress(children ...g.Node) g.NodeFunc {
	return g.El("progress", children...)
}

func Script(children ...g.Node) g.NodeFunc {
	return g.El("script", children...)
}

func Section(children ...g.Node) g.NodeFunc {
	return g.El("section", children...)
}

func Select(children ...g.Node) g.NodeFunc {
	return g.El("select", children...)
}

func Source(children ...g.Node) g.NodeFunc {
	return g.El("source", children...)
}

func Span(children ...g.Node) g.NodeFunc {
	return g.El("span", children...)
}

func StyleEl(children ...g.Node) g.NodeFunc {
	return g.El("style", children...)
}

func Summary(children ...g.Node) g.NodeFunc {
	return g.El("summary", children...)
}

func SVG(children ...g.Node) g.NodeFunc {
	return g.El("svg", children...)
}

func Table(children ...g.Node) g.NodeFunc {
	return g.El("table", children...)
}

func TBody(children ...g.Node) g.NodeFunc {
	return g.El("tbody", children...)
}

func Td(children ...g.Node) g.NodeFunc {
	return g.El("td", children...)
}

func Textarea(children ...g.Node) g.NodeFunc {
	return g.El("textarea", children...)
}

func TFoot(children ...g.Node) g.NodeFunc {
	return g.El("tfoot", children...)
}

func Th(children ...g.Node) g.NodeFunc {
	return g.El("th", children...)
}

func THead(children ...g.Node) g.NodeFunc {
	return g.El("thead", children...)
}

func Tr(children ...g.Node) g.NodeFunc {
	return g.El("tr", children...)
}

func Ul(children ...g.Node) g.NodeFunc {
	return g.El("ul", children...)
}

func Wbr(children ...g.Node) g.NodeFunc {
	return g.El("wbr", children...)
}

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

func H1(text string, children ...g.Node) g.NodeFunc {
	return g.El("h1", g.Text(text), g.Group(children))
}

func H2(text string, children ...g.Node) g.NodeFunc {
	return g.El("h2", g.Text(text), g.Group(children))
}

func H3(text string, children ...g.Node) g.NodeFunc {
	return g.El("h3", g.Text(text), g.Group(children))
}

func H4(text string, children ...g.Node) g.NodeFunc {
	return g.El("h4", g.Text(text), g.Group(children))
}

func H5(text string, children ...g.Node) g.NodeFunc {
	return g.El("h5", g.Text(text), g.Group(children))
}

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

func TitleEl(title string, children ...g.Node) g.NodeFunc {
	return g.El("title", g.Text(title), g.Group(children))
}

func U(text string, children ...g.Node) g.NodeFunc {
	return g.El("u", g.Text(text), g.Group(children))
}

func Var(text string, children ...g.Node) g.NodeFunc {
	return g.El("var", g.Text(text), g.Group(children))
}
