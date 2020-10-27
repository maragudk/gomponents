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

func TestBr(t *testing.T) {
	t.Run("returns a br element in context", func(t *testing.T) {
		assert.Equal(t, `<p>Test<br />Me</p>`, el.P(g.Text("Test"), el.Br(), g.Text("Me")))
	})
}

func TestHr(t *testing.T) {
	t.Run("returns a hr element with class", func(t *testing.T) {
		assert.Equal(t, `<hr class="test" />`, el.Hr(g.Attr("class", "test")))
	})
}

func TestBlockQuote(t *testing.T) {
	t.Run("returns a blockquote element", func(t *testing.T) {
		assert.Equal(t, `<blockquote>hat</blockquote>`, el.BlockQuote(g.Text("hat")))
	})
}

func TestDd(t *testing.T) {
	t.Run("returns a dd element", func(t *testing.T) {
		assert.Equal(t, `<dd>hat</dd>`, el.Dd("hat"))
	})
}

func TestDl(t *testing.T) {
	t.Run("returns a dl element", func(t *testing.T) {
		assert.Equal(t, `<dl><dt>hat</dt><dd>a nice thing for the head</dd></dl>`, el.Dl(el.Dt("hat"), el.Dd("a nice thing for the head")))
	})
}

func TestDt(t *testing.T) {
	t.Run("returns a dt element", func(t *testing.T) {
		assert.Equal(t, `<dt>hat</dt>`, el.Dt("hat"))
	})
}

func TestFigCaption(t *testing.T) {
	t.Run("returns a figcaption element", func(t *testing.T) {
		assert.Equal(t, `<figcaption>hat</figcaption>`, el.FigCaption(g.Text("hat")))
	})
}

func TestFigure(t *testing.T) {
	t.Run("returns a figure element", func(t *testing.T) {
		assert.Equal(t, `<figure>hat</figure>`, el.FigCaption(g.Text("hat")))
	})
}

func TestPre(t *testing.T) {
	t.Run("returns a pre element", func(t *testing.T) {
		assert.Equal(t, `<pre>hat</pre>`, el.Pre(g.Text("hat")))
	})
}
