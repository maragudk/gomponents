//go:build go1.18
// +build go1.18

package gomponents_test

import (
	"os"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/internal/assert"
)

func TestMap(t *testing.T) {
	t.Run("maps slices to nodes", func(t *testing.T) {
		items := []string{"hat", "partyhat", "turtlehat"}
		lis := g.Map(items, func(i string) g.Node {
			return g.El("li", g.Text(i))
		})

		list := g.El("ul", lis...)

		assert.Equal(t, `<ul><li>hat</li><li>partyhat</li><li>turtlehat</li></ul>`, list)
	})
}

func ExampleMap() {
	items := []string{"party hat", "super hat"}
	e := g.El("ul", g.Group(g.Map(items, func(i string) g.Node {
		return g.El("li", g.Text(i))
	})))
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>super hat</li></ul>
}
