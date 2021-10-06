package components_test

import (
	"os"
	"testing"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/internal/assert"
)

func TestClasses(t *testing.T) {
	t.Run("given a map, returns sorted keys from the map with value true", func(t *testing.T) {
		assert.Equal(t, ` class="boheme-hat hat partyhat"`, c.Classes{
			"boheme-hat": true,
			"hat":        true,
			"partyhat":   true,
			"turtlehat":  false,
		})
	})

	t.Run("renders as attribute in an element", func(t *testing.T) {
		e := g.El("div", c.Classes{"hat": true})
		assert.Equal(t, `<div class="hat"></div>`, e)
	})

	t.Run("also works with fmt", func(t *testing.T) {
		a := c.Classes{"hat": true}
		if a.String() != ` class="hat"` {
			t.FailNow()
		}
	})
}

func ExampleClasses() {
	e := g.El("div", c.Classes{"party-hat": true, "boring-hat": false})
	_ = e.Render(os.Stdout)
	// Output: <div class="party-hat"></div>
}
