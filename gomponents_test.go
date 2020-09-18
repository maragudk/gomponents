package gomponents_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
)

func TestAttr(t *testing.T) {
	t.Run("renders just the local name with one argument", func(t *testing.T) {
		a := g.Attr("required")
		assert.Equal(t, " required", a)
	})

	t.Run("renders the name and value when given two arguments", func(t *testing.T) {
		a := g.Attr("id", "hat")
		assert.Equal(t, ` id="hat"`, a)
	})

	t.Run("panics with more than two arguments", func(t *testing.T) {
		called := false
		defer func() {
			if err := recover(); err != nil {
				called = true
			}
		}()
		g.Attr("name", "value", "what is this?")
		if !called {
			t.FailNow()
		}
	})
}

func TestEl(t *testing.T) {
	t.Run("renders an empty element if no children given", func(t *testing.T) {
		e := g.El("div")
		assert.Equal(t, "<div/>", e)
	})

	t.Run("renders an empty element if only attributes given as children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"/>`, e)
	})

	t.Run("renders an element, attributes, and element children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"), g.El("span"))
		assert.Equal(t, `<div class="hat"><span/></div>`, e)
	})

	t.Run("renders attributes at the correct place regardless of placement in parameter list", func(t *testing.T) {
		e := g.El("div", g.El("span"), g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"><span/></div>`, e)
	})
}

func TestText(t *testing.T) {
	t.Run("renders escaped text", func(t *testing.T) {
		e := g.Text("<div/>")
		assert.Equal(t, "&lt;div/&gt;", e)
	})
}

func TestRaw(t *testing.T) {
	t.Run("renders raw text", func(t *testing.T) {
		e := g.Raw("<div/>")
		assert.Equal(t, "<div/>", e)
	})
}
