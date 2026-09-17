package html_test

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"maragu.dev/gomponents/internal/assert"
)

type erroringWriter struct{}

func (w *erroringWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("don't want to write")
}

func TestDoctype(t *testing.T) {
	t.Run("returns doctype and children", func(t *testing.T) {
		assert.Equal(t, `<!doctype html><html></html>`, Doctype(g.El("html")))
	})

	t.Run("ignores a nil sibling", func(t *testing.T) {
		assert.Equal(t, `<!doctype html>`, Doctype(nil))
		assert.Equal(t, `<!doctype html>`, Doctype(g.If(false, g.El("html"))))
	})

	t.Run("errors on write error in Render", func(t *testing.T) {
		err := Doctype(g.El("html")).Render(&erroringWriter{})
		assert.Error(t, err)
	})
}

func TestSimpleElements(t *testing.T) {
	tests := []struct {
		Name     string
		Func     func(...g.Node) g.Node
		Expected string
	}{
		{Name: "a", Func: A, Expected: `<a></a>`},
		{Name: "abbr", Func: Abbr, Expected: `<abbr></abbr>`},
		{Name: "address", Func: Address, Expected: `<address></address>`},
		{Name: "article", Func: Article, Expected: `<article></article>`},
		{Name: "aside", Func: Aside, Expected: `<aside></aside>`},
		{Name: "audio", Func: Audio, Expected: `<audio></audio>`},
		{Name: "b", Func: B, Expected: `<b></b>`},
		{Name: "blockquote", Func: BlockQuote, Expected: `<blockquote></blockquote>`},
		{Name: "body", Func: Body, Expected: `<body></body>`},
		{Name: "button", Func: Button, Expected: `<button></button>`},
		{Name: "canvas", Func: Canvas, Expected: `<canvas></canvas>`},
		{Name: "caption", Func: Caption, Expected: `<caption></caption>`},
		{Name: "cite", Func: Cite, Expected: `<cite></cite>`},
		{Name: "cite", Func: CiteEl, Expected: `<cite></cite>`},
		{Name: "code", Func: Code, Expected: `<code></code>`},
		{Name: "colgroup", Func: ColGroup, Expected: `<colgroup></colgroup>`},
		{Name: "data", Func: DataEl, Expected: `<data></data>`},
		{Name: "datalist", Func: DataList, Expected: `<datalist></datalist>`},
		{Name: "dd", Func: Dd, Expected: `<dd></dd>`},
		{Name: "del", Func: Del, Expected: `<del></del>`},
		{Name: "details", Func: Details, Expected: `<details></details>`},
		{Name: "dfn", Func: Dfn, Expected: `<dfn></dfn>`},
		{Name: "dialog", Func: Dialog, Expected: `<dialog></dialog>`},
		{Name: "div", Func: Div, Expected: `<div></div>`},
		{Name: "dl", Func: Dl, Expected: `<dl></dl>`},
		{Name: "dt", Func: Dt, Expected: `<dt></dt>`},
		{Name: "em", Func: Em, Expected: `<em></em>`},
		{Name: "fieldset", Func: FieldSet, Expected: `<fieldset></fieldset>`},
		{Name: "figcaption", Func: FigCaption, Expected: `<figcaption></figcaption>`},
		{Name: "figure", Func: Figure, Expected: `<figure></figure>`},
		{Name: "footer", Func: Footer, Expected: `<footer></footer>`},
		{Name: "form", Func: Form, Expected: `<form></form>`},
		{Name: "form", Func: FormEl, Expected: `<form></form>`},
		{Name: "h1", Func: H1, Expected: `<h1></h1>`},
		{Name: "h2", Func: H2, Expected: `<h2></h2>`},
		{Name: "h3", Func: H3, Expected: `<h3></h3>`},
		{Name: "h4", Func: H4, Expected: `<h4></h4>`},
		{Name: "h5", Func: H5, Expected: `<h5></h5>`},
		{Name: "h6", Func: H6, Expected: `<h6></h6>`},
		{Name: "head", Func: Head, Expected: `<head></head>`},
		{Name: "header", Func: Header, Expected: `<header></header>`},
		{Name: "hgroup", Func: HGroup, Expected: `<hgroup></hgroup>`},
		{Name: "html", Func: HTML, Expected: `<html></html>`},
		{Name: "i", Func: I, Expected: `<i></i>`},
		{Name: "iframe", Func: IFrame, Expected: `<iframe></iframe>`},
		{Name: "ins", Func: Ins, Expected: `<ins></ins>`},
		{Name: "kbd", Func: Kbd, Expected: `<kbd></kbd>`},
		{Name: "label", Func: Label, Expected: `<label></label>`},
		{Name: "label", Func: LabelEl, Expected: `<label></label>`},
		{Name: "legend", Func: Legend, Expected: `<legend></legend>`},
		{Name: "li", Func: Li, Expected: `<li></li>`},
		{Name: "main", Func: Main, Expected: `<main></main>`},
		{Name: "mark", Func: Mark, Expected: `<mark></mark>`},
		{Name: "menu", Func: Menu, Expected: `<menu></menu>`},
		{Name: "meter", Func: Meter, Expected: `<meter></meter>`},
		{Name: "nav", Func: Nav, Expected: `<nav></nav>`},
		{Name: "noscript", Func: NoScript, Expected: `<noscript></noscript>`},
		{Name: "object", Func: Object, Expected: `<object></object>`},
		{Name: "ol", Func: Ol, Expected: `<ol></ol>`},
		{Name: "optgroup", Func: OptGroup, Expected: `<optgroup></optgroup>`},
		{Name: "option", Func: Option, Expected: `<option></option>`},
		{Name: "output", Func: Output, Expected: `<output></output>`},
		{Name: "p", Func: P, Expected: `<p></p>`},
		{Name: "picture", Func: Picture, Expected: `<picture></picture>`},
		{Name: "pre", Func: Pre, Expected: `<pre></pre>`},
		{Name: "progress", Func: Progress, Expected: `<progress></progress>`},
		{Name: "q", Func: Q, Expected: `<q></q>`},
		{Name: "rp", Func: Rp, Expected: `<rp></rp>`},
		{Name: "rt", Func: Rt, Expected: `<rt></rt>`},
		{Name: "ruby", Func: Ruby, Expected: `<ruby></ruby>`},
		{Name: "s", Func: S, Expected: `<s></s>`},
		{Name: "samp", Func: Samp, Expected: `<samp></samp>`},
		{Name: "script", Func: Script, Expected: `<script></script>`},
		{Name: "search", Func: Search, Expected: `<search></search>`},
		{Name: "section", Func: Section, Expected: `<section></section>`},
		{Name: "select", Func: Select, Expected: `<select></select>`},
		{Name: "slot", Func: SlotEl, Expected: `<slot></slot>`},
		{Name: "small", Func: Small, Expected: `<small></small>`},
		{Name: "span", Func: Span, Expected: `<span></span>`},
		{Name: "strong", Func: Strong, Expected: `<strong></strong>`},
		{Name: "style", Func: StyleEl, Expected: `<style></style>`},
		{Name: "sub", Func: Sub, Expected: `<sub></sub>`},
		{Name: "summary", Func: Summary, Expected: `<summary></summary>`},
		{Name: "sup", Func: Sup, Expected: `<sup></sup>`},
		{Name: "svg", Func: SVG, Expected: `<svg></svg>`},
		{Name: "table", Func: Table, Expected: `<table></table>`},
		{Name: "tbody", Func: TBody, Expected: `<tbody></tbody>`},
		{Name: "td", Func: Td, Expected: `<td></td>`},
		{Name: "template", Func: Template, Expected: `<template></template>`},
		{Name: "textarea", Func: Textarea, Expected: `<textarea></textarea>`},
		{Name: "tfoot", Func: TFoot, Expected: `<tfoot></tfoot>`},
		{Name: "th", Func: Th, Expected: `<th></th>`},
		{Name: "thead", Func: THead, Expected: `<thead></thead>`},
		{Name: "time", Func: Time, Expected: `<time></time>`},
		{Name: "title", Func: TitleEl, Expected: `<title></title>`},
		{Name: "tr", Func: Tr, Expected: `<tr></tr>`},
		{Name: "u", Func: U, Expected: `<u></u>`},
		{Name: "ul", Func: Ul, Expected: `<ul></ul>`},
		{Name: "var", Func: Var, Expected: `<var></var>`},
		{Name: "video", Func: Video, Expected: `<video></video>`},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("renders with children", func(t *testing.T) {
				n := test.Func(g.Attr("id", "hat"))
				assert.Equal(t, fmt.Sprintf(`<%v id="hat"></%v>`, test.Name, test.Name), n)
			})

			t.Run("renders an empty element when no children show up", func(t *testing.T) {
				assert.Equal(t, test.Expected, test.Func())
			})

			t.Run("allocates nothing when no children show up", func(t *testing.T) {
				assertNoAllocs(t, func() {
					sink = test.Func()
				})
			})
		})
	}
}

func TestSimpleVoidKindElements(t *testing.T) {
	tests := []struct {
		Name     string
		Func     func(...g.Node) g.Node
		Expected string
	}{
		{Name: "area", Func: Area, Expected: `<area>`},
		{Name: "base", Func: Base, Expected: `<base>`},
		{Name: "br", Func: Br, Expected: `<br>`},
		{Name: "col", Func: Col, Expected: `<col>`},
		{Name: "embed", Func: Embed, Expected: `<embed>`},
		{Name: "hr", Func: Hr, Expected: `<hr>`},
		{Name: "img", Func: Img, Expected: `<img>`},
		{Name: "input", Func: Input, Expected: `<input>`},
		{Name: "link", Func: Link, Expected: `<link>`},
		{Name: "meta", Func: Meta, Expected: `<meta>`},
		{Name: "param", Func: Param, Expected: `<param>`},
		{Name: "source", Func: Source, Expected: `<source>`},
		{Name: "wbr", Func: Wbr, Expected: `<wbr>`},
		{Name: "track", Func: Track, Expected: `<track>`},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("renders with children", func(t *testing.T) {
				n := test.Func(g.Attr("id", "hat"))
				assert.Equal(t, fmt.Sprintf(`<%v id="hat">`, test.Name), n)
			})

			t.Run("renders the lone start tag when no children show up", func(t *testing.T) {
				assert.Equal(t, test.Expected, test.Func())
			})

			t.Run("allocates nothing when no children show up", func(t *testing.T) {
				assertNoAllocs(t, func() {
					sink = test.Func()
				})
			})
		})
	}
}

func TestSharedNodes(t *testing.T) {
	tests := []struct {
		Name     string
		Node     g.Node
		Expected string
	}{
		{Name: "a childless element", Node: Br(), Expected: `<br>`},
		{Name: "a boolean attribute", Node: Async(), Expected: ` async`},
	}

	for _, test := range tests {
		t.Run(test.Name+" renders the same into every writer that asks for it, having no state to carry between them", func(t *testing.T) {
			// Render into separate writers at the same time, so that the race detector
			// gets a say in whether a shared node is really stateless.
			results := make([]string, 8)

			var wg sync.WaitGroup
			for i := range results {
				wg.Add(1)

				go func(i int) {
					defer wg.Done()

					var b strings.Builder
					if err := test.Node.Render(&b); err != nil {
						t.Error("error rendering:", err)
						return
					}
					results[i] = b.String()
				}(i)
			}
			wg.Wait()

			for i, result := range results {
				if result != test.Expected {
					t.Fatalf(`expected writer %v to get "%v" but it got "%v"`, i, test.Expected, result)
				}
			}
		})
	}
}

// sink keeps a node returned inside [testing.AllocsPerRun] alive, so that the compiler
// cannot decide the call has no effect and drop it, making the allocation count a lie.
var sink g.Node

// assertNoAllocs fails if calling f allocates anything at all.
func assertNoAllocs(t *testing.T, f func()) {
	t.Helper()

	if allocs := testing.AllocsPerRun(100, f); allocs != 0 {
		t.Fatalf("expected 0 allocations but got %v", allocs)
	}
}
