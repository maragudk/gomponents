// Package html provides common HTML elements and attributes.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element for a list of elements.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes for a list of attributes.
package html

import (
	"io"

	g "maragu.dev/gomponents"
)

// Doctype returns a special kind of [g.Node] that prefixes its sibling with the string "<!doctype html>".
func Doctype(sibling g.Node) g.Node {
	return g.NodeFunc(func(w io.Writer) error {
		if _, err := w.Write([]byte("<!doctype html>")); err != nil {
			return err
		}
		return sibling.Render(w)
	})
}

func A(children ...g.Node) g.Node {
	return g.El("a", children...)
}

func Address(children ...g.Node) g.Node {
	return g.El("address", children...)
}

func Area(children ...g.Node) g.Node {
	return g.El("area", children...)
}

func Article(children ...g.Node) g.Node {
	return g.El("article", children...)
}

func Aside(children ...g.Node) g.Node {
	return g.El("aside", children...)
}

func Audio(children ...g.Node) g.Node {
	return g.El("audio", children...)
}

func Base(children ...g.Node) g.Node {
	return g.El("base", children...)
}

func BlockQuote(children ...g.Node) g.Node {
	return g.El("blockquote", children...)
}

func Body(children ...g.Node) g.Node {
	return g.El("body", children...)
}

func Br(children ...g.Node) g.Node {
	return g.El("br", children...)
}

func Button(children ...g.Node) g.Node {
	return g.El("button", children...)
}

func Canvas(children ...g.Node) g.Node {
	return g.El("canvas", children...)
}

func Cite(children ...g.Node) g.Node {
	return g.El("cite", children...)
}

// Deprecated: Use [Cite] instead.
func CiteEl(children ...g.Node) g.Node {
	return Cite(children...)
}

func Code(children ...g.Node) g.Node {
	return g.El("code", children...)
}

func Col(children ...g.Node) g.Node {
	return g.El("col", children...)
}

func ColGroup(children ...g.Node) g.Node {
	return g.El("colgroup", children...)
}

func DataEl(children ...g.Node) g.Node {
	return g.El("data", children...)
}

func DataList(children ...g.Node) g.Node {
	return g.El("datalist", children...)
}

func Details(children ...g.Node) g.Node {
	return g.El("details", children...)
}

func Dialog(children ...g.Node) g.Node {
	return g.El("dialog", children...)
}

func Div(children ...g.Node) g.Node {
	return g.El("div", children...)
}

func Dl(children ...g.Node) g.Node {
	return g.El("dl", children...)
}

func Embed(children ...g.Node) g.Node {
	return g.El("embed", children...)
}

func Form(children ...g.Node) g.Node {
	return g.El("form", children...)
}

// Deprecated: Use [Form] instead.
func FormEl(children ...g.Node) g.Node {
	return Form(children...)
}

func FieldSet(children ...g.Node) g.Node {
	return g.El("fieldset", children...)
}

func Figure(children ...g.Node) g.Node {
	return g.El("figure", children...)
}

func Footer(children ...g.Node) g.Node {
	return g.El("footer", children...)
}

func Head(children ...g.Node) g.Node {
	return g.El("head", children...)
}

func Header(children ...g.Node) g.Node {
	return g.El("header", children...)
}

func HGroup(children ...g.Node) g.Node {
	return g.El("hgroup", children...)
}

func Hr(children ...g.Node) g.Node {
	return g.El("hr", children...)
}

func HTML(children ...g.Node) g.Node {
	return g.El("html", children...)
}

func IFrame(children ...g.Node) g.Node {
	return g.El("iframe", children...)
}

func Img(children ...g.Node) g.Node {
	return g.El("img", children...)
}

func Input(children ...g.Node) g.Node {
	return g.El("input", children...)
}

func Label(children ...g.Node) g.Node {
	return g.El("label", children...)
}

// Deprecated: Use [Label] instead.
func LabelEl(children ...g.Node) g.Node {
	return Label(children...)
}

func Legend(children ...g.Node) g.Node {
	return g.El("legend", children...)
}

func Li(children ...g.Node) g.Node {
	return g.El("li", children...)
}

func Link(children ...g.Node) g.Node {
	return g.El("link", children...)
}

func Main(children ...g.Node) g.Node {
	return g.El("main", children...)
}

func Menu(children ...g.Node) g.Node {
	return g.El("menu", children...)
}

func Meta(children ...g.Node) g.Node {
	return g.El("meta", children...)
}

func Meter(children ...g.Node) g.Node {
	return g.El("meter", children...)
}

func Nav(children ...g.Node) g.Node {
	return g.El("nav", children...)
}

func NoScript(children ...g.Node) g.Node {
	return g.El("noscript", children...)
}

func Object(children ...g.Node) g.Node {
	return g.El("object", children...)
}

func Ol(children ...g.Node) g.Node {
	return g.El("ol", children...)
}

func OptGroup(children ...g.Node) g.Node {
	return g.El("optgroup", children...)
}

func Option(children ...g.Node) g.Node {
	return g.El("option", children...)
}

func P(children ...g.Node) g.Node {
	return g.El("p", children...)
}

func Param(children ...g.Node) g.Node {
	return g.El("param", children...)
}

func Picture(children ...g.Node) g.Node {
	return g.El("picture", children...)
}

func Pre(children ...g.Node) g.Node {
	return g.El("pre", children...)
}

func Progress(children ...g.Node) g.Node {
	return g.El("progress", children...)
}

func Script(children ...g.Node) g.Node {
	return g.El("script", children...)
}

func Section(children ...g.Node) g.Node {
	return g.El("section", children...)
}

func Select(children ...g.Node) g.Node {
	return g.El("select", children...)
}

func SlotEl(children ...g.Node) g.Node {
	return g.El("slot", children...)
}

func Source(children ...g.Node) g.Node {
	return g.El("source", children...)
}

func Span(children ...g.Node) g.Node {
	return g.El("span", children...)
}

func StyleEl(children ...g.Node) g.Node {
	return g.El("style", children...)
}

func Summary(children ...g.Node) g.Node {
	return g.El("summary", children...)
}

func SVG(children ...g.Node) g.Node {
	return g.El("svg", children...)
}

func Table(children ...g.Node) g.Node {
	return g.El("table", children...)
}

func TBody(children ...g.Node) g.Node {
	return g.El("tbody", children...)
}

func Td(children ...g.Node) g.Node {
	return g.El("td", children...)
}

func Template(children ...g.Node) g.Node {
	return g.El("template", children...)
}

func Textarea(children ...g.Node) g.Node {
	return g.El("textarea", children...)
}

func TFoot(children ...g.Node) g.Node {
	return g.El("tfoot", children...)
}

func Th(children ...g.Node) g.Node {
	return g.El("th", children...)
}

func THead(children ...g.Node) g.Node {
	return g.El("thead", children...)
}

func Tr(children ...g.Node) g.Node {
	return g.El("tr", children...)
}

func Ul(children ...g.Node) g.Node {
	return g.El("ul", children...)
}

func Wbr(children ...g.Node) g.Node {
	return g.El("wbr", children...)
}

func Abbr(children ...g.Node) g.Node {
	return g.El("abbr", g.Group(children))
}

func B(children ...g.Node) g.Node {
	return g.El("b", g.Group(children))
}

func Caption(children ...g.Node) g.Node {
	return g.El("caption", g.Group(children))
}

func Dd(children ...g.Node) g.Node {
	return g.El("dd", g.Group(children))
}

func Del(children ...g.Node) g.Node {
	return g.El("del", g.Group(children))
}

func Dfn(children ...g.Node) g.Node {
	return g.El("dfn", g.Group(children))
}

func Dt(children ...g.Node) g.Node {
	return g.El("dt", g.Group(children))
}

func Em(children ...g.Node) g.Node {
	return g.El("em", g.Group(children))
}

func FigCaption(children ...g.Node) g.Node {
	return g.El("figcaption", g.Group(children))
}

func H1(children ...g.Node) g.Node {
	return g.El("h1", g.Group(children))
}

func H2(children ...g.Node) g.Node {
	return g.El("h2", g.Group(children))
}

func H3(children ...g.Node) g.Node {
	return g.El("h3", g.Group(children))
}

func H4(children ...g.Node) g.Node {
	return g.El("h4", g.Group(children))
}

func H5(children ...g.Node) g.Node {
	return g.El("h5", g.Group(children))
}

func H6(children ...g.Node) g.Node {
	return g.El("h6", g.Group(children))
}

func I(children ...g.Node) g.Node {
	return g.El("i", g.Group(children))
}

func Ins(children ...g.Node) g.Node {
	return g.El("ins", g.Group(children))
}

func Kbd(children ...g.Node) g.Node {
	return g.El("kbd", g.Group(children))
}

func Mark(children ...g.Node) g.Node {
	return g.El("mark", g.Group(children))
}

func Q(children ...g.Node) g.Node {
	return g.El("q", g.Group(children))
}

func S(children ...g.Node) g.Node {
	return g.El("s", g.Group(children))
}

func Samp(children ...g.Node) g.Node {
	return g.El("samp", g.Group(children))
}

func Small(children ...g.Node) g.Node {
	return g.El("small", g.Group(children))
}

func Strong(children ...g.Node) g.Node {
	return g.El("strong", g.Group(children))
}

func Sub(children ...g.Node) g.Node {
	return g.El("sub", g.Group(children))
}

func Sup(children ...g.Node) g.Node {
	return g.El("sup", g.Group(children))
}

func Time(children ...g.Node) g.Node {
	return g.El("time", g.Group(children))
}

func TitleEl(children ...g.Node) g.Node {
	return g.El("title", g.Group(children))
}

func U(children ...g.Node) g.Node {
	return g.El("u", g.Group(children))
}

func Var(children ...g.Node) g.Node {
	return g.El("var", g.Group(children))
}

func Video(children ...g.Node) g.Node {
	return g.El("video", g.Group(children))
}
