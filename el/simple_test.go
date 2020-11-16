package el_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestSimpleElements(t *testing.T) {
	cases := map[string]func(...g.Node) g.NodeFunc{
		"address":    el.Address,
		"article":    el.Article,
		"aside":      el.Aside,
		"audio":      el.Audio,
		"blockquote": el.BlockQuote,
		"body":       el.Body,
		"button":     el.Button,
		"canvas":     el.Canvas,
		"cite":       el.Cite,
		"code":       el.Code,
		"colgroup":   el.ColGroup,
		"data":       el.Data,
		"datalist":   el.DataList,
		"details":    el.Details,
		"dialog":     el.Dialog,
		"div":        el.Div,
		"dl":         el.Dl,
		"fieldset":   el.FieldSet,
		"figure":     el.Figure,
		"footer":     el.Footer,
		"head":       el.Head,
		"header":     el.Header,
		"hgroup":     el.HGroup,
		"html":       el.HTML,
		"iframe":     el.IFrame,
		"legend":     el.Legend,
		"li":         el.Li,
		"main":       el.Main,
		"menu":       el.Menu,
		"meter":      el.Meter,
		"nav":        el.Nav,
		"noscript":   el.NoScript,
		"object":     el.Object,
		"ol":         el.Ol,
		"optgroup":   el.OptGroup,
		"p":          el.P,
		"picture":    el.Picture,
		"pre":        el.Pre,
		"script":     el.Script,
		"section":    el.Section,
		"span":       el.Span,
		"style":      el.Style,
		"summary":    el.Summary,
		"table":      el.Table,
		"tbody":      el.TBody,
		"td":         el.Td,
		"tfoot":      el.TFoot,
		"th":         el.Th,
		"thead":      el.THead,
		"tr":         el.Tr,
		"ul":         el.Ul,
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
		"area":   el.Area,
		"base":   el.Base,
		"br":     el.Br,
		"col":    el.Col,
		"embed":  el.Embed,
		"hr":     el.Hr,
		"link":   el.Link,
		"meta":   el.Meta,
		"param":  el.Param,
		"source": el.Source,
		"wbr":    el.Wbr,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat">`, name), n)
		})
	}
}
