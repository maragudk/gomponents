package components_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	c "github.com/maragudk/gomponents/components"
)

func TestInputHidden(t *testing.T) {
	t.Run("returns an input element with type hidden, and the given name and value", func(t *testing.T) {
		n := c.InputHidden("id", "partyhat", g.Attr("class", "hat"))
		assert.Equal(t, `<input type="hidden" name="id" value="partyhat" class="hat">`, n)
	})
}
