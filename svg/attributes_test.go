package svg_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/internal/assert"
	. "github.com/maragudk/gomponents/svg"
)

func TestSimpleAttributes(t *testing.T) {
	cases := map[string]func(string) g.Node{
		"clip-rule":    ClipRule,
		"d":            D,
		"fill":         Fill,
		"fill-rule":    FillRule,
		"stroke":       Stroke,
		"stroke-width": StrokeWidth,
		"viewBox":      ViewBox,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf(`should output %v="hat"`, name), func(t *testing.T) {
			n := g.El("element", fn("hat"))
			assert.Equal(t, fmt.Sprintf(`<element %v="hat"></element>`, name), n)
		})
	}
}
