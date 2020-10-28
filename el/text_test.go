package el_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestTextElements(t *testing.T) {
	cases := map[string]func(string, ...g.Node) g.NodeFunc{
		"abbr":       el.Abbr,
		"b":          el.B,
		"caption":    el.Caption,
		"dd":         el.Dd,
		"del":        el.Del,
		"dfn":        el.Dfn,
		"dt":         el.Dt,
		"em":         el.Em,
		"figcaption": el.FigCaption,
		"h1":         el.H1,
		"h2":         el.H2,
		"h3":         el.H3,
		"h4":         el.H4,
		"h5":         el.H5,
		"h6":         el.H6,
		"i":          el.I,
		"ins":        el.Ins,
		"kbd":        el.Kbd,
		"mark":       el.Mark,
		"q":          el.Q,
		"s":          el.S,
		"samp":       el.Samp,
		"small":      el.Small,
		"strong":     el.Strong,
		"sub":        el.Sub,
		"sup":        el.Sup,
		"time":       el.Time,
		"title":      el.Title,
		"u":          el.U,
		"var":        el.Var,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn("hat", g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat">hat</%v>`, name, name), n)
		})
	}
}
