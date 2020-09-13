package gomponents_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
)

func TestAttr(t *testing.T) {
	t.Run("renders just the local name with one argument", func(t *testing.T) {
		a := g.Attr("required")
		equal(t, " required", a.Render())
	})

	t.Run("renders the name and value when given two arguments", func(t *testing.T) {
		a := g.Attr("id", "hat")
		equal(t, ` id="hat"`, a.Render())
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
		equal(t, "<div/>", e.Render())
	})

	t.Run("renders an empty element if only attributes given as children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"))
		equal(t, `<div class="hat"/>`, e.Render())
	})

	t.Run("renders an element, attributes, and element children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"), g.El("span"))
		equal(t, `<div class="hat"><span/></div>`, e.Render())
	})

	t.Run("renders attributes at the correct place regardless of placement in parameter list", func(t *testing.T) {
		e := g.El("div", g.El("span"), g.Attr("class", "hat"))
		equal(t, `<div class="hat"><span/></div>`, e.Render())
	})
}

func TestText(t *testing.T) {
	t.Run("renders escaped text", func(t *testing.T) {
		e := g.Text("<div/>")
		equal(t, "&lt;div/&gt;", e.Render())
	})
}

func TestRaw(t *testing.T) {
	t.Run("renders raw text", func(t *testing.T) {
		e := g.Raw("<div/>")
		equal(t, "<div/>", e.Render())
	})
}

func equal(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
		t.FailNow()
	}
}
