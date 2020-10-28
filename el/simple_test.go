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
		"area":       el.Area,
		"article":    el.Article,
		"aside":      el.Aside,
		"audio":      el.Audio,
		"base":       el.Base,
		"blockquote": el.BlockQuote,
		"body":       el.Body,
		"br":         el.Br,
		"button":     el.Button,
		"canvas":     el.Canvas,
		"cite":       el.Cite,
		"code":       el.Code,
		"col":        el.Col,
		"colgroup":   el.ColGroup,
		"data":       el.Data,
		"datalist":   el.DataList,
		"details":    el.Details,
		"dialog":     el.Dialog,
		"div":        el.Div,
		"dl":         el.Dl,
		"embed":      el.Embed,
		"fieldset":   el.FieldSet,
		"figure":     el.Figure,
		"footer":     el.Footer,
		"head":       el.Head,
		"header":     el.Header,
		"hgroup":     el.HGroup,
		"hr":         el.Hr,
		"html":       el.HTML,
		"iframe":     el.IFrame,
		"legend":     el.Legend,
		"li":         el.Li,
		"link":       el.Link,
		"main":       el.Main,
		"menu":       el.Menu,
		"meta":       el.Meta,
		"meter":      el.Meter,
		"nav":        el.Nav,
		"noscript":   el.NoScript,
		"object":     el.Object,
		"ol":         el.Ol,
		"optgroup":   el.OptGroup,
		"p":          el.P,
		"param":      el.Param,
		"picture":    el.Picture,
		"pre":        el.Pre,
		"script":     el.Script,
		"section":    el.Section,
		"source":     el.Source,
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
		"wbr":        el.Wbr,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat" />`, name), n)
		})
	}
}
