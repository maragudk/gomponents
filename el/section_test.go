package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

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
