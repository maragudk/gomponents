package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestSpan(t *testing.T) {
	t.Run("returns a span element", func(t *testing.T) {
		assert.Equal(t, `<span>hat</span>`, el.Span(g.Text("hat")))
	})
}

func TestA(t *testing.T) {
	t.Run("returns an a element with a href attribute", func(t *testing.T) {
		assert.Equal(t, `<a href="#">hat</a>`, el.A("#", g.Text("hat")))
	})
}

func TestB(t *testing.T) {
	t.Run("returns a b element", func(t *testing.T) {
		assert.Equal(t, `<b id="text">hat</b>`, el.B("hat", g.Attr("id", "text")))
	})
}

func TestStrong(t *testing.T) {
	t.Run("returns a strong element", func(t *testing.T) {
		assert.Equal(t, `<strong id="text">hat</strong>`, el.Strong("hat", g.Attr("id", "text")))
	})
}

func TestI(t *testing.T) {
	t.Run("returns an i element", func(t *testing.T) {
		assert.Equal(t, `<i id="text">hat</i>`, el.I("hat", g.Attr("id", "text")))
	})
}

func TestEm(t *testing.T) {
	t.Run("returns an em element", func(t *testing.T) {
		assert.Equal(t, `<em id="text">hat</em>`, el.Em("hat", g.Attr("id", "text")))
	})
}
