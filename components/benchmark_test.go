//go:build go1.24

package components_test

import (
	"bufio"
	"io"
	"strconv"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func BenchmarkJoinAttrs(b *testing.B) {
	flat := func() []g.Node {
		return []g.Node{
			Class("inline-flex items-center"),
			ID("badge"),
			Class("px-2.5 py-0.5 rounded-full"),
			Type("button"),
			Class("bg-green-100 text-green-800"),
		}
	}

	nested := func() []g.Node {
		return []g.Node{
			Class("inline-flex items-center"),
			g.Group{
				ID("badge"),
				g.Group{
					Class("px-2.5 py-0.5 rounded-full"),
					Class("bg-green-100 text-green-800"),
				},
			},
			Type("button"),
		}
	}

	noMatch := func() []g.Node {
		return []g.Node{ID("badge"), Type("button"), Lang("en"), Title("hat")}
	}

	cases := []struct {
		Name     string
		Children func() []g.Node
	}{
		{Name: "flat", Children: flat},
		{Name: "nested groups", Children: nested},
		{Name: "no match", Children: noMatch},
	}

	for _, c := range cases {
		b.Run(c.Name+"/construct", func(b *testing.B) {
			children := c.Children()

			for b.Loop() {
				_ = JoinAttrs("class", children...)
			}
		})

		b.Run(c.Name+"/construct and render", func(b *testing.B) {
			children := c.Children()

			for b.Loop() {
				_ = JoinAttrs("class", children...).Render(io.Discard)
			}
		})
	}
}

// staticTree is a medium-sized tree with nothing dynamic in it: a document head and a navigation
// with a couple of dozen links, the kind of thing [Static] is for. It is built anew on each call,
// like a component called per request.
func staticTree() g.Node {
	links := make([]g.Node, 24)
	for i := range links {
		links[i] = Li(A(Href("/docs/section-"+strconv.Itoa(i)), Class("nav-link"), g.Text("Section "+strconv.Itoa(i))))
	}
	return g.Group{
		Head(
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			TitleEl(g.Text("Party hats")),
			Link(Rel("stylesheet"), Href("/static/app.css")),
			Link(Rel("icon"), Href("/static/favicon.ico")),
			Script(Src("/static/app.js"), Defer()),
		),
		Nav(Class("navbar"), Ul(Class("nav"), g.Group(links))),
	}
}

func BenchmarkStatic(b *testing.B) {
	// The buffered writer makes the many short writes of a direct render cost something, unlike
	// [io.Discard], and is the size of the buffer net/http puts in front of a response writer.
	writers := []struct {
		Name string
		New  func() io.Writer
	}{
		{Name: "discarded", New: func() io.Writer { return io.Discard }},
		{Name: "buffered", New: func() io.Writer { return bufio.NewWriterSize(io.Discard, 2048) }},
	}

	for _, w := range writers {
		// Built once and rendered repeatedly, to separate the cost of rendering from the cost
		// of building the tree, which Static does not skip.
		b.Run("direct/render pre-built/"+w.Name, func(b *testing.B) {
			tree := staticTree()
			w := w.New()

			for b.Loop() {
				_ = tree.Render(w)
			}
		})

		b.Run("static/render pre-built/"+w.Name, func(b *testing.B) {
			var slot string
			node := Static(&slot, staticTree())
			w := w.New()

			for b.Loop() {
				_ = node.Render(w)
			}
		})

		// Built and rendered on every iteration, like a component called per request.
		b.Run("direct/construct and render/"+w.Name, func(b *testing.B) {
			w := w.New()

			for b.Loop() {
				_ = staticTree().Render(w)
			}
		})

		b.Run("static/construct and render/"+w.Name, func(b *testing.B) {
			var slot string
			w := w.New()

			for b.Loop() {
				_ = Static(&slot, staticTree()).Render(w)
			}
		})
	}
}
