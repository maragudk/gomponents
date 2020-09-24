package attr_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/attr"
)

func TestForms(t *testing.T) {
	t.Run("adds placeholder and required attributes", func(t *testing.T) {
		e := g.El("input", attr.Placeholder("hat"), attr.Required())
		assert.Equal(t, `<input placeholder="hat" required/>`, e)
	})
}
