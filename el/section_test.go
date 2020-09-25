package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestAddress(t *testing.T) {
	t.Run("returns an address element", func(t *testing.T) {
		assert.Equal(t, `<address />`, el.Address())
	})
}

func TestArticle(t *testing.T) {
	t.Run("returns an article element", func(t *testing.T) {
		assert.Equal(t, `<article />`, el.Article())
	})
}

func TestAside(t *testing.T) {
	t.Run("returns an aside element", func(t *testing.T) {
		assert.Equal(t, `<aside />`, el.Aside())
	})
}

func TestFooter(t *testing.T) {
	t.Run("returns a footer element", func(t *testing.T) {
		assert.Equal(t, `<footer />`, el.Footer())
	})
}

func TestHeader(t *testing.T) {
	t.Run("returns a header element", func(t *testing.T) {
		assert.Equal(t, `<header />`, el.Header())
	})
}

func TestH1(t *testing.T) {
	t.Run("returns an h1 element", func(t *testing.T) {
		assert.Equal(t, `<h1 id="headline">hat</h1>`, el.H1("hat", g.Attr("id", "headline")))
	})
}

func TestH2(t *testing.T) {
	t.Run("returns an h2 element", func(t *testing.T) {
		assert.Equal(t, `<h2 id="headline">hat</h2>`, el.H2("hat", g.Attr("id", "headline")))
	})
}

func TestH3(t *testing.T) {
	t.Run("returns an h3 element", func(t *testing.T) {
		assert.Equal(t, `<h3 id="headline">hat</h3>`, el.H3("hat", g.Attr("id", "headline")))
	})
}

func TestH4(t *testing.T) {
	t.Run("returns an h4 element", func(t *testing.T) {
		assert.Equal(t, `<h4 id="headline">hat</h4>`, el.H4("hat", g.Attr("id", "headline")))
	})
}

func TestH5(t *testing.T) {
	t.Run("returns an h5 element", func(t *testing.T) {
		assert.Equal(t, `<h5 id="headline">hat</h5>`, el.H5("hat", g.Attr("id", "headline")))
	})
}

func TestH6(t *testing.T) {
	t.Run("returns an h6 element", func(t *testing.T) {
		assert.Equal(t, `<h6 id="headline">hat</h6>`, el.H6("hat", g.Attr("id", "headline")))
	})
}

func TestHGroup(t *testing.T) {
	t.Run("returns an hgroup element", func(t *testing.T) {
		assert.Equal(t, `<hgroup />`, el.HGroup())
	})
}

func TestMainEl(t *testing.T) {
	t.Run("returns a main element", func(t *testing.T) {
		assert.Equal(t, `<main />`, el.Main())
	})
}

func TestNav(t *testing.T) {
	t.Run("returns a nav element", func(t *testing.T) {
		assert.Equal(t, `<nav />`, el.Nav())
	})
}

func TestSection(t *testing.T) {
	t.Run("returns a section element", func(t *testing.T) {
		assert.Equal(t, `<section />`, el.Section())
	})
}
