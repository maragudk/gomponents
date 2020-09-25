package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestDocument(t *testing.T) {
	t.Run("returns doctype and children", func(t *testing.T) {
		assert.Equal(t, `<!doctype html><html />`, el.Document(g.El("html")))
	})
}

func TestHTML(t *testing.T) {
	t.Run("returns an html element", func(t *testing.T) {
		assert.Equal(t, "<html><div /><span /></html>", el.HTML(g.El("div"), g.El("span")))
	})
}

func TestHead(t *testing.T) {
	t.Run("returns a head element", func(t *testing.T) {
		assert.Equal(t, "<head><title /><link /></head>", el.Head(g.El("title"), g.El("link")))
	})
}

func TestBody(t *testing.T) {
	t.Run("returns a body element", func(t *testing.T) {
		assert.Equal(t, "<body><div /><span /></body>", el.Body(g.El("div"), g.El("span")))
	})
}

func TestTitle(t *testing.T) {
	t.Run("returns a title element with text content", func(t *testing.T) {
		assert.Equal(t, "<title>hat</title>", el.Title("hat"))
	})
}

func TestMeta(t *testing.T) {
	t.Run("returns a meta element", func(t *testing.T) {
		assert.Equal(t, `<meta charset="utf-8" />`, el.Meta(g.Attr("charset", "utf-8")))
	})
}

func TestLink(t *testing.T) {
	t.Run("returns a link element", func(t *testing.T) {
		assert.Equal(t, `<link rel="stylesheet" />`, el.Link(g.Attr("rel", "stylesheet")))
	})
}

func TestStyle(t *testing.T) {
	t.Run("returns a style element", func(t *testing.T) {
		assert.Equal(t, `<style type="text/css" />`, el.Style(g.Attr("type", "text/css")))
	})
}

func TestBase(t *testing.T) {
	t.Run("returns a base element", func(t *testing.T) {
		assert.Equal(t, `<base href="/hat/" />`, el.Base(g.Attr("href", "/hat/")))
	})
}
