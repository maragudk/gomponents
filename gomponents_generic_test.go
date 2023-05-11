//go:build go1.18
// +build go1.18

package gomponents_test

import (
	"fmt"
	"os"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/internal/assert"
)

func TestMap(t *testing.T) {
	t.Run("maps slices to nodes", func(t *testing.T) {
		items := []string{"hat", "partyhat", "turtlehat"}
		lis := g.Map(items, func(item string) g.Node {
			return g.El("li", g.Text(item))
		})

		list := g.El("ul", lis...)

		assert.Equal(t, `<ul><li>hat</li><li>partyhat</li><li>turtlehat</li></ul>`, list)
	})
}

func ExampleMap() {
	items := []string{"party hat", "super hat"}
	e := g.El("ul", g.Group(g.Map(items, func(item string) g.Node {
		return g.El("li", g.Text(item))
	})))
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>super hat</li></ul>
}

func TestMapIndex(t *testing.T) {
	t.Run("maps slices to nodes", func(t *testing.T) {
		items := []string{"hat", "partyhat", "turtlehat"}
		lis := g.MapIndex(items, func(item string, i int) g.Node {
			return g.El("li", g.Attr("data-index", fmt.Sprint(i)), g.Text(item))
		})

		list := g.El("ul", lis...)

		assert.Equal(t, `<ul><li data-index="0">hat</li><li data-index="1">partyhat</li><li data-index="2">turtlehat</li></ul>`, list)
	})
}

func ExampleMapIndex() {
	items := []string{"party hat", "super hat"}
	e := g.El("ul", g.Group(g.MapIndex(items, func(item string, i int) g.Node {
		return g.El("li",
			g.If(i%2 == 0, g.Attr("class", "even")),
			g.If(i%2 == 1, g.Attr("class", "odd")),
			g.Text(item),
		)
	})))
	_ = e.Render(os.Stdout)
	// Output: <ul><li class="even">party hat</li><li class="odd">super hat</li></ul>
}
