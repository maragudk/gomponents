package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestDocument(t *testing.T) {
	t.Run("returns doctype and children", func(t *testing.T) {
		assert.Equal(t, `<!doctype html><html/>`, el.Document(g.El("html")))
	})
}

func TestHTML(t *testing.T) {
	t.Run("returns an html element", func(t *testing.T) {
		assert.Equal(t, "<html><div/><span/></html>", el.HTML(g.El("div"), g.El("span")))
	})
}

func TestHead(t *testing.T) {
	t.Run("returns a head element", func(t *testing.T) {
		assert.Equal(t, "<head><title/><link/></head>", el.Head(g.El("title"), g.El("link")))
	})
}

func TestBody(t *testing.T) {
	t.Run("returns a body element", func(t *testing.T) {
		assert.Equal(t, "<body><div/><span/></body>", el.Body(g.El("div"), g.El("span")))
	})
}

func TestTitle(t *testing.T) {
	t.Run("returns a title element with text content", func(t *testing.T) {
		assert.Equal(t, "<title>hat</title>", el.Title("hat"))
	})
}

func TestMeta(t *testing.T) {
	t.Run("returns a meta element", func(t *testing.T) {
		assert.Equal(t, `<meta charset="utf-8"/>`, el.Meta(g.Attr("charset", "utf-8")))
	})
}

func TestLink(t *testing.T) {
	t.Run("returns a link element", func(t *testing.T) {
		assert.Equal(t, `<link rel="stylesheet"/>`, el.Link(g.Attr("rel", "stylesheet")))
	})
}

func TestStyle(t *testing.T) {
	t.Run("returns a style element", func(t *testing.T) {
		assert.Equal(t, `<style type="text/css"/>`, el.Style(g.Attr("type", "text/css")))
	})
}

func TestDiv(t *testing.T) {
	t.Run("returns a div element", func(t *testing.T) {
		assert.Equal(t, `<div><span/></div>`, el.Div(el.Span()))
	})
}

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

func TestP(t *testing.T) {
	t.Run("returns a p element", func(t *testing.T) {
		assert.Equal(t, `<p>hat</p>`, el.P(g.Text("hat")))
	})
}

func TestH1(t *testing.T) {
	t.Run("returns an h1 element", func(t *testing.T) {
		assert.Equal(t, `<h1>hat</h1>`, el.H1("hat"))
	})
}

func TestH2(t *testing.T) {
	t.Run("returns an h2 element", func(t *testing.T) {
		assert.Equal(t, `<h2>hat</h2>`, el.H2("hat"))
	})
}

func TestH3(t *testing.T) {
	t.Run("returns an h3 element", func(t *testing.T) {
		assert.Equal(t, `<h3>hat</h3>`, el.H3("hat"))
	})
}

func TestH4(t *testing.T) {
	t.Run("returns an h4 element", func(t *testing.T) {
		assert.Equal(t, `<h4>hat</h4>`, el.H4("hat"))
	})
}

func TestH5(t *testing.T) {
	t.Run("returns an h5 element", func(t *testing.T) {
		assert.Equal(t, `<h5>hat</h5>`, el.H5("hat"))
	})
}

func TestH6(t *testing.T) {
	t.Run("returns an h6 element", func(t *testing.T) {
		assert.Equal(t, `<h6>hat</h6>`, el.H6("hat"))
	})
}

func TestOl(t *testing.T) {
	t.Run("returns an ol element", func(t *testing.T) {
		assert.Equal(t, `<ol><li/></ol>`, el.Ol(el.Li()))
	})
}

func TestUl(t *testing.T) {
	t.Run("returns a ul element", func(t *testing.T) {
		assert.Equal(t, `<ul><li/></ul>`, el.Ul(el.Li()))
	})
}

func TestLi(t *testing.T) {
	t.Run("returns an li element", func(t *testing.T) {
		assert.Equal(t, `<li>hat</li>`, el.Li(g.Text("hat")))
	})
}

func TestB(t *testing.T) {
	t.Run("returns a b element", func(t *testing.T) {
		assert.Equal(t, `<b>hat</b>`, el.B("hat"))
	})
}

func TestStrong(t *testing.T) {
	t.Run("returns a strong element", func(t *testing.T) {
		assert.Equal(t, `<strong>hat</strong>`, el.Strong("hat"))
	})
}

func TestI(t *testing.T) {
	t.Run("returns an i element", func(t *testing.T) {
		assert.Equal(t, `<i>hat</i>`, el.I("hat"))
	})
}

func TestEm(t *testing.T) {
	t.Run("returns an em element", func(t *testing.T) {
		assert.Equal(t, `<em>hat</em>`, el.Em("hat"))
	})
}

func TestImg(t *testing.T) {
	t.Run("returns an img element with href and alt attributes", func(t *testing.T) {
		assert.Equal(t, `<img src="hat.png" alt="hat"/>`, el.Img("hat.png", "hat"))
	})
}
