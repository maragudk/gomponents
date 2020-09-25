package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestDiv(t *testing.T) {
	t.Run("returns a div element", func(t *testing.T) {
		assert.Equal(t, `<div><span /></div>`, el.Div(el.Span()))
	})
}

func TestOl(t *testing.T) {
	t.Run("returns an ol element", func(t *testing.T) {
		assert.Equal(t, `<ol><li /></ol>`, el.Ol(el.Li()))
	})
}

func TestUl(t *testing.T) {
	t.Run("returns a ul element", func(t *testing.T) {
		assert.Equal(t, `<ul><li /></ul>`, el.Ul(el.Li()))
	})
}

func TestLi(t *testing.T) {
	t.Run("returns an li element", func(t *testing.T) {
		assert.Equal(t, `<li>hat</li>`, el.Li(g.Text("hat")))
	})
}

func TestP(t *testing.T) {
	t.Run("returns a p element", func(t *testing.T) {
		assert.Equal(t, `<p>hat</p>`, el.P(g.Text("hat")))
	})
}
