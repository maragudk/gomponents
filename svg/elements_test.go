package svg_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/internal/assert"
	. "github.com/maragudk/gomponents/svg"
)

func TestSimpleElements(t *testing.T) {
	cases := map[string]func(...g.Node) g.Node{
		"path": Path,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := fn(g.Attr("id", "hat"))
			assert.Equal(t, fmt.Sprintf(`<%v id="hat"></%v>`, name, name), n)
		})
	}
}

func TestSVG(t *testing.T) {
	t.Run("outputs svg element with xml namespace attribute", func(t *testing.T) {
		assert.Equal(t, `<svg xmlns="http://www.w3.org/2000/svg"><path></path></svg>`, SVG(g.El("path")))
	})
}
