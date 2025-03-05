package html_test

import (
	"errors"
	"fmt"
	"strings"
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

	t.Run("errors on write error in Render", func(t *testing.T) {
		err := Doctype(g.El("html")).Render(&erroringWriter{})
		assert.Error(t, err)
	})
}

func TestSimpleElements(t *testing.T) {
	tests := []struct {
		Name string
		Func func(...g.Node) g.Node
	}{
		{Name: "a", Func: A},
		{Name: "abbr", Func: Abbr},
		{Name: "address", Func: Address},
		{Name: "article", Func: Article},
		{Name: "aside", Func: Aside},
		{Name: "audio", Func: Audio},
		{Name: "b", Func: B},
		{Name: "blockquote", Func: BlockQuote},
		{Name: "body", Func: Body},
		{Name: "button", Func: Button},
		{Name: "canvas", Func: Canvas},
		{Name: "caption", Func: Caption},
		{Name: "cite", Func: Cite},
		{Name: "cite", Func: CiteEl},
		{Name: "code", Func: Code},
		{Name: "colgroup", Func: ColGroup},
		{Name: "data", Func: DataEl},
		{Name: "datalist", Func: DataList},
		{Name: "dd", Func: Dd},
		{Name: "del", Func: Del},
		{Name: "details", Func: Details},
		{Name: "dfn", Func: Dfn},
		{Name: "dialog", Func: Dialog},
		{Name: "div", Func: Div},
		{Name: "dl", Func: Dl},
		{Name: "dt", Func: Dt},
		{Name: "em", Func: Em},
		{Name: "fieldset", Func: FieldSet},
		{Name: "figcaption", Func: FigCaption},
		{Name: "figure", Func: Figure},
		{Name: "footer", Func: Footer},
		{Name: "form", Func: Form},
		{Name: "form", Func: FormEl},
		{Name: "h1", Func: H1},
		{Name: "h2", Func: H2},
		{Name: "h3", Func: H3},
		{Name: "h4", Func: H4},
		{Name: "h5", Func: H5},
		{Name: "h6", Func: H6},
		{Name: "head", Func: Head},
		{Name: "header", Func: Header},
		{Name: "hgroup", Func: HGroup},
		{Name: "html", Func: HTML},
		{Name: "i", Func: I},
		{Name: "iframe", Func: IFrame},
		{Name: "ins", Func: Ins},
		{Name: "kbd", Func: Kbd},
		{Name: "label", Func: Label},
		{Name: "label", Func: LabelEl},
		{Name: "legend", Func: Legend},
		{Name: "li", Func: Li},
		{Name: "main", Func: Main},
		{Name: "mark", Func: Mark},
		{Name: "menu", Func: Menu},
		{Name: "meter", Func: Meter},
		{Name: "nav", Func: Nav},
		{Name: "noscript", Func: NoScript},
		{Name: "object", Func: Object},
		{Name: "ol", Func: Ol},
		{Name: "optgroup", Func: OptGroup},
		{Name: "option", Func: Option},
		{Name: "p", Func: P},
		{Name: "picture", Func: Picture},
		{Name: "pre", Func: Pre},
		{Name: "progress", Func: Progress},
		{Name: "q", Func: Q},
		{Name: "s", Func: S},
		{Name: "samp", Func: Samp},
		{Name: "script", Func: Script},
		{Name: "section", Func: Section},
		{Name: "select", Func: Select},
		{Name: "slot", Func: SlotEl},
		{Name: "small", Func: Small},
		{Name: "span", Func: Span},
		{Name: "strong", Func: Strong},
		{Name: "style", Func: StyleEl},
		{Name: "sub", Func: Sub},
		{Name: "summary", Func: Summary},
		{Name: "sup", Func: Sup},
		{Name: "svg", Func: SVG},
		{Name: "table", Func: Table},
		{Name: "tbody", Func: TBody},
		{Name: "td", Func: Td},
		{Name: "template", Func: Template},
		{Name: "textarea", Func: Textarea},
		{Name: "tfoot", Func: TFoot},
		{Name: "th", Func: Th},
		{Name: "thead", Func: THead},
		{Name: "time", Func: Time},
		{Name: "title", Func: TitleEl},
		{Name: "tr", Func: Tr},
		{Name: "u", Func: U},
		{Name: "ul", Func: Ul},
		{Name: "var", Func: Var},
		{Name: "video", Func: Video},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			n := test.Func(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat"></%v>`, test.Name, test.Name), n)
		})
	}
}

func TestSimpleVoidKindElements(t *testing.T) {
	tests := []struct {
		Name string
		Func func(...g.Node) g.Node
	}{
		{Name: "area", Func: Area},
		{Name: "base", Func: Base},
		{Name: "br", Func: Br},
		{Name: "col", Func: Col},
		{Name: "embed", Func: Embed},
		{Name: "hr", Func: Hr},
		{Name: "img", Func: Img},
		{Name: "input", Func: Input},
		{Name: "link", Func: Link},
		{Name: "meta", Func: Meta},
		{Name: "param", Func: Param},
		{Name: "source", Func: Source},
		{Name: "wbr", Func: Wbr},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			n := test.Func(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat">`, test.Name), n)
		})
	}
}

func BenchmarkLargeHTMLDocument(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		elements := make([]g.Node, 0, 10000)

		for i := 0; i < 5000; i++ {
			elements = append(elements,
				Div(Class("foo")),
				Span(Class("bar")),
			)
		}
		doc := Div(elements...)

		var sb strings.Builder
		_ = doc.Render(&sb)
	}
}
