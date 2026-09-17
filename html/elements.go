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
// If sibling is nil, only the doctype is rendered.
func Doctype(sibling g.Node) g.Node {
	return g.NodeFunc(func(w io.Writer) error {
		if _, err := io.WriteString(w, "<!doctype html>"); err != nil {
			return err
		}

		if sibling == nil {
			return nil
		}

		return sibling.Render(w)
	})
}

func A(children ...g.Node) g.Node {
	if len(children) == 0 {
		return aEl
	}
	return g.El("a", children...)
}

var aEl = g.El("a")

func Address(children ...g.Node) g.Node {
	if len(children) == 0 {
		return addressEl
	}
	return g.El("address", children...)
}

var addressEl = g.El("address")

func Area(children ...g.Node) g.Node {
	if len(children) == 0 {
		return areaEl
	}
	return g.El("area", children...)
}

var areaEl = g.El("area")

func Article(children ...g.Node) g.Node {
	if len(children) == 0 {
		return articleEl
	}
	return g.El("article", children...)
}

var articleEl = g.El("article")

func Aside(children ...g.Node) g.Node {
	if len(children) == 0 {
		return asideEl
	}
	return g.El("aside", children...)
}

var asideEl = g.El("aside")

func Audio(children ...g.Node) g.Node {
	if len(children) == 0 {
		return audioEl
	}
	return g.El("audio", children...)
}

var audioEl = g.El("audio")

func Base(children ...g.Node) g.Node {
	if len(children) == 0 {
		return baseEl
	}
	return g.El("base", children...)
}

var baseEl = g.El("base")

func BlockQuote(children ...g.Node) g.Node {
	if len(children) == 0 {
		return blockquoteEl
	}
	return g.El("blockquote", children...)
}

var blockquoteEl = g.El("blockquote")

func Body(children ...g.Node) g.Node {
	if len(children) == 0 {
		return bodyEl
	}
	return g.El("body", children...)
}

var bodyEl = g.El("body")

func Br(children ...g.Node) g.Node {
	if len(children) == 0 {
		return brEl
	}
	return g.El("br", children...)
}

var brEl = g.El("br")

func Button(children ...g.Node) g.Node {
	if len(children) == 0 {
		return buttonEl
	}
	return g.El("button", children...)
}

var buttonEl = g.El("button")

func Canvas(children ...g.Node) g.Node {
	if len(children) == 0 {
		return canvasEl
	}
	return g.El("canvas", children...)
}

var canvasEl = g.El("canvas")

func Cite(children ...g.Node) g.Node {
	if len(children) == 0 {
		return citeEl
	}
	return g.El("cite", children...)
}

var citeEl = g.El("cite")

// Deprecated: Use [Cite] instead.
//
//go:fix inline
func CiteEl(children ...g.Node) g.Node {
	return Cite(children...)
}

func Code(children ...g.Node) g.Node {
	if len(children) == 0 {
		return codeEl
	}
	return g.El("code", children...)
}

var codeEl = g.El("code")

func Col(children ...g.Node) g.Node {
	if len(children) == 0 {
		return colEl
	}
	return g.El("col", children...)
}

var colEl = g.El("col")

func ColGroup(children ...g.Node) g.Node {
	if len(children) == 0 {
		return colgroupEl
	}
	return g.El("colgroup", children...)
}

var colgroupEl = g.El("colgroup")

func DataEl(children ...g.Node) g.Node {
	if len(children) == 0 {
		return dataEl
	}
	return g.El("data", children...)
}

var dataEl = g.El("data")

func DataList(children ...g.Node) g.Node {
	if len(children) == 0 {
		return datalistEl
	}
	return g.El("datalist", children...)
}

var datalistEl = g.El("datalist")

func Details(children ...g.Node) g.Node {
	if len(children) == 0 {
		return detailsEl
	}
	return g.El("details", children...)
}

var detailsEl = g.El("details")

func Dialog(children ...g.Node) g.Node {
	if len(children) == 0 {
		return dialogEl
	}
	return g.El("dialog", children...)
}

var dialogEl = g.El("dialog")

func Div(children ...g.Node) g.Node {
	if len(children) == 0 {
		return divEl
	}
	return g.El("div", children...)
}

var divEl = g.El("div")

func Dl(children ...g.Node) g.Node {
	if len(children) == 0 {
		return dlEl
	}
	return g.El("dl", children...)
}

var dlEl = g.El("dl")

func Embed(children ...g.Node) g.Node {
	if len(children) == 0 {
		return embedEl
	}
	return g.El("embed", children...)
}

var embedEl = g.El("embed")

func Form(children ...g.Node) g.Node {
	if len(children) == 0 {
		return formEl
	}
	return g.El("form", children...)
}

var formEl = g.El("form")

// Deprecated: Use [Form] instead.
//
//go:fix inline
func FormEl(children ...g.Node) g.Node {
	return Form(children...)
}

func FieldSet(children ...g.Node) g.Node {
	if len(children) == 0 {
		return fieldsetEl
	}
	return g.El("fieldset", children...)
}

var fieldsetEl = g.El("fieldset")

func Figure(children ...g.Node) g.Node {
	if len(children) == 0 {
		return figureEl
	}
	return g.El("figure", children...)
}

var figureEl = g.El("figure")

func Footer(children ...g.Node) g.Node {
	if len(children) == 0 {
		return footerEl
	}
	return g.El("footer", children...)
}

var footerEl = g.El("footer")

func Head(children ...g.Node) g.Node {
	if len(children) == 0 {
		return headEl
	}
	return g.El("head", children...)
}

var headEl = g.El("head")

func Header(children ...g.Node) g.Node {
	if len(children) == 0 {
		return headerEl
	}
	return g.El("header", children...)
}

var headerEl = g.El("header")

func HGroup(children ...g.Node) g.Node {
	if len(children) == 0 {
		return hgroupEl
	}
	return g.El("hgroup", children...)
}

var hgroupEl = g.El("hgroup")

func Hr(children ...g.Node) g.Node {
	if len(children) == 0 {
		return hrEl
	}
	return g.El("hr", children...)
}

var hrEl = g.El("hr")

func HTML(children ...g.Node) g.Node {
	if len(children) == 0 {
		return htmlEl
	}
	return g.El("html", children...)
}

var htmlEl = g.El("html")

func IFrame(children ...g.Node) g.Node {
	if len(children) == 0 {
		return iframeEl
	}
	return g.El("iframe", children...)
}

var iframeEl = g.El("iframe")

func Img(children ...g.Node) g.Node {
	if len(children) == 0 {
		return imgEl
	}
	return g.El("img", children...)
}

var imgEl = g.El("img")

func Input(children ...g.Node) g.Node {
	if len(children) == 0 {
		return inputEl
	}
	return g.El("input", children...)
}

var inputEl = g.El("input")

func Label(children ...g.Node) g.Node {
	if len(children) == 0 {
		return labelEl
	}
	return g.El("label", children...)
}

var labelEl = g.El("label")

// Deprecated: Use [Label] instead.
//
//go:fix inline
func LabelEl(children ...g.Node) g.Node {
	return Label(children...)
}

func Legend(children ...g.Node) g.Node {
	if len(children) == 0 {
		return legendEl
	}
	return g.El("legend", children...)
}

var legendEl = g.El("legend")

func Li(children ...g.Node) g.Node {
	if len(children) == 0 {
		return liEl
	}
	return g.El("li", children...)
}

var liEl = g.El("li")

func Link(children ...g.Node) g.Node {
	if len(children) == 0 {
		return linkEl
	}
	return g.El("link", children...)
}

var linkEl = g.El("link")

func Main(children ...g.Node) g.Node {
	if len(children) == 0 {
		return mainEl
	}
	return g.El("main", children...)
}

var mainEl = g.El("main")

func Menu(children ...g.Node) g.Node {
	if len(children) == 0 {
		return menuEl
	}
	return g.El("menu", children...)
}

var menuEl = g.El("menu")

func Meta(children ...g.Node) g.Node {
	if len(children) == 0 {
		return metaEl
	}
	return g.El("meta", children...)
}

var metaEl = g.El("meta")

func Meter(children ...g.Node) g.Node {
	if len(children) == 0 {
		return meterEl
	}
	return g.El("meter", children...)
}

var meterEl = g.El("meter")

func Nav(children ...g.Node) g.Node {
	if len(children) == 0 {
		return navEl
	}
	return g.El("nav", children...)
}

var navEl = g.El("nav")

func NoScript(children ...g.Node) g.Node {
	if len(children) == 0 {
		return noscriptEl
	}
	return g.El("noscript", children...)
}

var noscriptEl = g.El("noscript")

func Object(children ...g.Node) g.Node {
	if len(children) == 0 {
		return objectEl
	}
	return g.El("object", children...)
}

var objectEl = g.El("object")

func Ol(children ...g.Node) g.Node {
	if len(children) == 0 {
		return olEl
	}
	return g.El("ol", children...)
}

var olEl = g.El("ol")

func OptGroup(children ...g.Node) g.Node {
	if len(children) == 0 {
		return optgroupEl
	}
	return g.El("optgroup", children...)
}

var optgroupEl = g.El("optgroup")

func Output(children ...g.Node) g.Node {
	if len(children) == 0 {
		return outputEl
	}
	return g.El("output", children...)
}

var outputEl = g.El("output")

func Option(children ...g.Node) g.Node {
	if len(children) == 0 {
		return optionEl
	}
	return g.El("option", children...)
}

var optionEl = g.El("option")

func P(children ...g.Node) g.Node {
	if len(children) == 0 {
		return pEl
	}
	return g.El("p", children...)
}

var pEl = g.El("p")

func Param(children ...g.Node) g.Node {
	if len(children) == 0 {
		return paramEl
	}
	return g.El("param", children...)
}

var paramEl = g.El("param")

func Picture(children ...g.Node) g.Node {
	if len(children) == 0 {
		return pictureEl
	}
	return g.El("picture", children...)
}

var pictureEl = g.El("picture")

func Pre(children ...g.Node) g.Node {
	if len(children) == 0 {
		return preEl
	}
	return g.El("pre", children...)
}

var preEl = g.El("pre")

func Progress(children ...g.Node) g.Node {
	if len(children) == 0 {
		return progressEl
	}
	return g.El("progress", children...)
}

var progressEl = g.El("progress")

func Rp(children ...g.Node) g.Node {
	if len(children) == 0 {
		return rpEl
	}
	return g.El("rp", children...)
}

var rpEl = g.El("rp")

func Rt(children ...g.Node) g.Node {
	if len(children) == 0 {
		return rtEl
	}
	return g.El("rt", children...)
}

var rtEl = g.El("rt")

func Ruby(children ...g.Node) g.Node {
	if len(children) == 0 {
		return rubyEl
	}
	return g.El("ruby", children...)
}

var rubyEl = g.El("ruby")

func Script(children ...g.Node) g.Node {
	if len(children) == 0 {
		return scriptEl
	}
	return g.El("script", children...)
}

var scriptEl = g.El("script")

func Search(children ...g.Node) g.Node {
	if len(children) == 0 {
		return searchEl
	}
	return g.El("search", children...)
}

var searchEl = g.El("search")

func Section(children ...g.Node) g.Node {
	if len(children) == 0 {
		return sectionEl
	}
	return g.El("section", children...)
}

var sectionEl = g.El("section")

func Select(children ...g.Node) g.Node {
	if len(children) == 0 {
		return selectEl
	}
	return g.El("select", children...)
}

var selectEl = g.El("select")

func SlotEl(children ...g.Node) g.Node {
	if len(children) == 0 {
		return slotEl
	}
	return g.El("slot", children...)
}

var slotEl = g.El("slot")

func Source(children ...g.Node) g.Node {
	if len(children) == 0 {
		return sourceEl
	}
	return g.El("source", children...)
}

var sourceEl = g.El("source")

func Span(children ...g.Node) g.Node {
	if len(children) == 0 {
		return spanEl
	}
	return g.El("span", children...)
}

var spanEl = g.El("span")

func StyleEl(children ...g.Node) g.Node {
	if len(children) == 0 {
		return styleEl
	}
	return g.El("style", children...)
}

var styleEl = g.El("style")

func Summary(children ...g.Node) g.Node {
	if len(children) == 0 {
		return summaryEl
	}
	return g.El("summary", children...)
}

var summaryEl = g.El("summary")

func SVG(children ...g.Node) g.Node {
	if len(children) == 0 {
		return svgEl
	}
	return g.El("svg", children...)
}

var svgEl = g.El("svg")

func Table(children ...g.Node) g.Node {
	if len(children) == 0 {
		return tableEl
	}
	return g.El("table", children...)
}

var tableEl = g.El("table")

func TBody(children ...g.Node) g.Node {
	if len(children) == 0 {
		return tbodyEl
	}
	return g.El("tbody", children...)
}

var tbodyEl = g.El("tbody")

func Td(children ...g.Node) g.Node {
	if len(children) == 0 {
		return tdEl
	}
	return g.El("td", children...)
}

var tdEl = g.El("td")

func Template(children ...g.Node) g.Node {
	if len(children) == 0 {
		return templateEl
	}
	return g.El("template", children...)
}

var templateEl = g.El("template")

func Textarea(children ...g.Node) g.Node {
	if len(children) == 0 {
		return textareaEl
	}
	return g.El("textarea", children...)
}

var textareaEl = g.El("textarea")

func TFoot(children ...g.Node) g.Node {
	if len(children) == 0 {
		return tfootEl
	}
	return g.El("tfoot", children...)
}

var tfootEl = g.El("tfoot")

func Th(children ...g.Node) g.Node {
	if len(children) == 0 {
		return thEl
	}
	return g.El("th", children...)
}

var thEl = g.El("th")

func THead(children ...g.Node) g.Node {
	if len(children) == 0 {
		return theadEl
	}
	return g.El("thead", children...)
}

var theadEl = g.El("thead")

func Tr(children ...g.Node) g.Node {
	if len(children) == 0 {
		return trEl
	}
	return g.El("tr", children...)
}

var trEl = g.El("tr")

func Ul(children ...g.Node) g.Node {
	if len(children) == 0 {
		return ulEl
	}
	return g.El("ul", children...)
}

var ulEl = g.El("ul")

func Wbr(children ...g.Node) g.Node {
	if len(children) == 0 {
		return wbrEl
	}
	return g.El("wbr", children...)
}

var wbrEl = g.El("wbr")

func Abbr(children ...g.Node) g.Node {
	if len(children) == 0 {
		return abbrEl
	}
	return g.El("abbr", g.Group(children))
}

var abbrEl = g.El("abbr")

func B(children ...g.Node) g.Node {
	if len(children) == 0 {
		return bEl
	}
	return g.El("b", g.Group(children))
}

var bEl = g.El("b")

func Caption(children ...g.Node) g.Node {
	if len(children) == 0 {
		return captionEl
	}
	return g.El("caption", g.Group(children))
}

var captionEl = g.El("caption")

func Dd(children ...g.Node) g.Node {
	if len(children) == 0 {
		return ddEl
	}
	return g.El("dd", g.Group(children))
}

var ddEl = g.El("dd")

func Del(children ...g.Node) g.Node {
	if len(children) == 0 {
		return delEl
	}
	return g.El("del", g.Group(children))
}

var delEl = g.El("del")

func Dfn(children ...g.Node) g.Node {
	if len(children) == 0 {
		return dfnEl
	}
	return g.El("dfn", g.Group(children))
}

var dfnEl = g.El("dfn")

func Dt(children ...g.Node) g.Node {
	if len(children) == 0 {
		return dtEl
	}
	return g.El("dt", g.Group(children))
}

var dtEl = g.El("dt")

func Em(children ...g.Node) g.Node {
	if len(children) == 0 {
		return emEl
	}
	return g.El("em", g.Group(children))
}

var emEl = g.El("em")

func FigCaption(children ...g.Node) g.Node {
	if len(children) == 0 {
		return figcaptionEl
	}
	return g.El("figcaption", g.Group(children))
}

var figcaptionEl = g.El("figcaption")

func H1(children ...g.Node) g.Node {
	if len(children) == 0 {
		return h1El
	}
	return g.El("h1", g.Group(children))
}

var h1El = g.El("h1")

func H2(children ...g.Node) g.Node {
	if len(children) == 0 {
		return h2El
	}
	return g.El("h2", g.Group(children))
}

var h2El = g.El("h2")

func H3(children ...g.Node) g.Node {
	if len(children) == 0 {
		return h3El
	}
	return g.El("h3", g.Group(children))
}

var h3El = g.El("h3")

func H4(children ...g.Node) g.Node {
	if len(children) == 0 {
		return h4El
	}
	return g.El("h4", g.Group(children))
}

var h4El = g.El("h4")

func H5(children ...g.Node) g.Node {
	if len(children) == 0 {
		return h5El
	}
	return g.El("h5", g.Group(children))
}

var h5El = g.El("h5")

func H6(children ...g.Node) g.Node {
	if len(children) == 0 {
		return h6El
	}
	return g.El("h6", g.Group(children))
}

var h6El = g.El("h6")

func I(children ...g.Node) g.Node {
	if len(children) == 0 {
		return iEl
	}
	return g.El("i", g.Group(children))
}

var iEl = g.El("i")

func Ins(children ...g.Node) g.Node {
	if len(children) == 0 {
		return insEl
	}
	return g.El("ins", g.Group(children))
}

var insEl = g.El("ins")

func Kbd(children ...g.Node) g.Node {
	if len(children) == 0 {
		return kbdEl
	}
	return g.El("kbd", g.Group(children))
}

var kbdEl = g.El("kbd")

func Mark(children ...g.Node) g.Node {
	if len(children) == 0 {
		return markEl
	}
	return g.El("mark", g.Group(children))
}

var markEl = g.El("mark")

func Q(children ...g.Node) g.Node {
	if len(children) == 0 {
		return qEl
	}
	return g.El("q", g.Group(children))
}

var qEl = g.El("q")

func S(children ...g.Node) g.Node {
	if len(children) == 0 {
		return sEl
	}
	return g.El("s", g.Group(children))
}

var sEl = g.El("s")

func Samp(children ...g.Node) g.Node {
	if len(children) == 0 {
		return sampEl
	}
	return g.El("samp", g.Group(children))
}

var sampEl = g.El("samp")

func Small(children ...g.Node) g.Node {
	if len(children) == 0 {
		return smallEl
	}
	return g.El("small", g.Group(children))
}

var smallEl = g.El("small")

func Strong(children ...g.Node) g.Node {
	if len(children) == 0 {
		return strongEl
	}
	return g.El("strong", g.Group(children))
}

var strongEl = g.El("strong")

func Sub(children ...g.Node) g.Node {
	if len(children) == 0 {
		return subEl
	}
	return g.El("sub", g.Group(children))
}

var subEl = g.El("sub")

func Sup(children ...g.Node) g.Node {
	if len(children) == 0 {
		return supEl
	}
	return g.El("sup", g.Group(children))
}

var supEl = g.El("sup")

func Time(children ...g.Node) g.Node {
	if len(children) == 0 {
		return timeEl
	}
	return g.El("time", g.Group(children))
}

var timeEl = g.El("time")

func TitleEl(children ...g.Node) g.Node {
	if len(children) == 0 {
		return titleEl
	}
	return g.El("title", g.Group(children))
}

var titleEl = g.El("title")

func U(children ...g.Node) g.Node {
	if len(children) == 0 {
		return uEl
	}
	return g.El("u", g.Group(children))
}

var uEl = g.El("u")

func Var(children ...g.Node) g.Node {
	if len(children) == 0 {
		return varEl
	}
	return g.El("var", g.Group(children))
}

var varEl = g.El("var")

func Video(children ...g.Node) g.Node {
	if len(children) == 0 {
		return videoEl
	}
	return g.El("video", g.Group(children))
}

var videoEl = g.El("video")

func Track(children ...g.Node) g.Node {
	if len(children) == 0 {
		return trackEl
	}
	return g.El("track", children...)
}

var trackEl = g.El("track")
