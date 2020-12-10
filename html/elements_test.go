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

func TestFormEl(t *testing.T) {
	t.Run("returns a form element with action and method attributes", func(t *testing.T) {
		assert.Equal(t, `<form action="/" method="post"></form>`, FormEl("/", "post"))
	})
}

func TestInput(t *testing.T) {
	t.Run("returns an input element with attributes type and name", func(t *testing.T) {
		assert.Equal(t, `<input type="text" name="hat">`, Input("text", "hat"))
	})
}

func TestLabel(t *testing.T) {
	t.Run("returns a label element with attribute for", func(t *testing.T) {
		assert.Equal(t, `<label for="hat">Hat</label>`, Label("hat", g.Text("Hat")))
	})
}

func TestOption(t *testing.T) {
	t.Run("returns an option element with attribute label and content", func(t *testing.T) {
		assert.Equal(t, `<option value="hat">Hat</option>`, Option("Hat", "hat"))
	})
}

func TestProgress(t *testing.T) {
	t.Run("returns a progress element with attributes value and max", func(t *testing.T) {
		assert.Equal(t, `<progress value="5.5" max="10"></progress>`, Progress(5.5, 10))
	})
}

func TestSelect(t *testing.T) {
	t.Run("returns a select element with attribute name", func(t *testing.T) {
		assert.Equal(t, `<select name="hat"><option value="partyhat">Partyhat</option></select>`,
			Select("hat", Option("Partyhat", "partyhat")))
	})
}

func TestTextarea(t *testing.T) {
	t.Run("returns a textarea element with attribute name", func(t *testing.T) {
		assert.Equal(t, `<textarea name="hat"></textarea>`, Textarea("hat"))
	})
}

func TestA(t *testing.T) {
	t.Run("returns an a element with a href attribute", func(t *testing.T) {
		assert.Equal(t, `<a href="#">hat</a>`, A("#", g.Text("hat")))
	})
}

func TestImg(t *testing.T) {
	t.Run("returns an img element with href and alt attributes", func(t *testing.T) {
		assert.Equal(t, `<img src="hat.png" alt="hat" id="image">`, Img("hat.png", "hat", g.Attr("id", "image")))
	})
}

func TestSimpleElements(t *testing.T) {
	cases := map[string]func(...g.Node) g.NodeFunc{
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
		"data":       Data,
		"datalist":   DataList,
		"details":    Details,
		"dialog":     Dialog,
		"div":        Div,
		"dl":         Dl,
		"fieldset":   FieldSet,
		"figure":     Figure,
		"footer":     Footer,
		"head":       Head,
		"header":     Header,
		"hgroup":     HGroup,
		"html":       HTML,
		"iframe":     IFrame,
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
		"p":          P,
		"picture":    Picture,
		"pre":        Pre,
		"script":     Script,
		"section":    Section,
		"span":       Span,
		"style":      StyleEl,
		"summary":    Summary,
		"svg":        SVG,
		"table":      Table,
		"tbody":      TBody,
		"td":         Td,
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
	cases := map[string]func(...g.Node) g.NodeFunc{
		"area":   Area,
		"base":   Base,
		"br":     Br,
		"col":    Col,
		"embed":  Embed,
		"hr":     Hr,
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
	cases := map[string]func(string, ...g.Node) g.NodeFunc{
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
