package el_test

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
	"testing"
)

func TestBr(t *testing.T) {
	t.Run("returns a br element", func(t *testing.T) {
		assert.Equal(t, `<p>Test<br />Me</p>`, el.P(g.Text("Test"), el.Br(), g.Text("Me")))
	})
}

func TestHr(t *testing.T) {
	t.Run("returns a hr element", func(t *testing.T) {
		assert.Equal(t, `<hr class="test" />`, el.Hr(attr.Class("test")))
	})
}
