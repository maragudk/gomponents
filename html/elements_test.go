package html_test

import (
	"errors"
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	. "github.com/maragudk/gomponents/html"
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
	cases := map[string]func(...g.Node) g.Node{
		"a":          A,
		"address":    Address,
		"article":    Article,
		"aside":      Aside,
		"audio":      Audio,
		"blockquote": BlockQuote,
		"body":       Body,
		"button":     Button,
		"canvas":     Canvas,
		"cite":       Cite,
		"code":       Code,
		"colgroup":   ColGroup,
		"data":       DataEl,
		"datalist":   DataList,
		"details":    Details,
		"dialog":     Dialog,
		"div":        Div,
		"dl":         Dl,
		"fieldset":   FieldSet,
		"figure":     Figure,
		"footer":     Footer,
		"form":       FormEl,
		"head":       Head,
		"header":     Header,
		"hgroup":     HGroup,
		"html":       HTML,
		"iframe":     IFrame,
		"label":      Label,
		"legend":     Legend,
		"li":         Li,
		"main":       Main,
		"menu":       Menu,
		"meter":      Meter,
		"nav":        Nav,
		"noscript":   NoScript,
		"object":     Object,
		"ol":         Ol,
		"optgroup":   OptGroup,
		"option":     Option,
		"p":          P,
		"picture":    Picture,
		"pre":        Pre,
		"progress":   Progress,
		"script":     Script,
		"section":    Section,
		"select":     Select,
		"span":       Span,
		"style":      StyleEl,
		"summary":    Summary,
		"svg":        SVG,
		"table":      Table,
		"tbody":      TBody,
		"td":         Td,
		"textarea":   Textarea,
		"tfoot":      TFoot,
		"th":         Th,
		"thead":      THead,
		"tr":         Tr,
		"ul":         Ul,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat"></%v>`, name, name), n)
		})
	}
}

func TestSimpleVoidKindElements(t *testing.T) {
	cases := map[string]func(...g.Node) g.Node{
		"area":   Area,
		"base":   Base,
		"br":     Br,
		"col":    Col,
		"embed":  Embed,
		"hr":     Hr,
		"img":    Img,
		"input":  Input,
		"link":   Link,
		"meta":   Meta,
		"param":  Param,
		"source": Source,
		"wbr":    Wbr,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat">`, name), n)
		})
	}
}

func TestTextElements(t *testing.T) {
	cases := map[string]func(string, ...g.Node) g.Node{
		"abbr":       Abbr,
		"b":          B,
		"caption":    Caption,
		"dd":         Dd,
		"del":        Del,
		"dfn":        Dfn,
		"dt":         Dt,
		"em":         Em,
		"figcaption": FigCaption,
		"h1":         H1,
		"h2":         H2,
		"h3":         H3,
		"h4":         H4,
		"h5":         H5,
		"h6":         H6,
		"i":          I,
		"ins":        Ins,
		"kbd":        Kbd,
		"mark":       Mark,
		"q":          Q,
		"s":          S,
		"samp":       Samp,
		"small":      Small,
		"strong":     Strong,
		"sub":        Sub,
		"sup":        Sup,
		"time":       Time,
		"title":      TitleEl,
		"u":          U,
		"var":        Var,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn("hat", g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat">hat</%v>`, name, name), n)
		})
	}
}
