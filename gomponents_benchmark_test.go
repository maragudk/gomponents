//go:build go1.24

package gomponents_test

import (
	"bufio"
	"io"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

// writeOnly is an [io.Writer] without a WriteString method, like a struct that embeds an
// [io.Writer] and promotes only Write.
type writeOnly struct {
	w io.Writer
}

func (w writeOnly) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

// node keeps a constructed [g.Node] reachable across benchmark iterations.
var node g.Node

func BenchmarkAttr(b *testing.B) {
	// A boolean attribute has no value; the rest are name-value attributes with values
	// named after what they stress.
	values := []struct {
		Name    string
		Boolean bool
		Value   string
	}{
		{Name: "boolean", Boolean: true},
		{Name: "no escaping", Value: "party"},
		{Name: "long value needing no escaping", Value: strings.Repeat("a title with no quotes or apostrophes in it ", 4)},
		{Name: "needing escaping", Value: `"party" & fun`},
		{Name: "many characters needing escaping", Value: strings.Repeat(`"hat" & `, 6)},
		{Name: "long value with few characters needing escaping", Value: strings.Repeat("It's a title with quotes & apostrophes in it. ", 4)},
	}

	// The buffered writers are the size of the [bufio.Writer] that net/http puts in front
	// of a handler's response writer, its bufferBeforeChunkingSize.
	writers := []struct {
		Name string
		W    io.Writer
	}{
		{Name: "discarded", W: io.Discard},
		{Name: "buffered", W: bufio.NewWriterSize(io.Discard, 2048)},
		{Name: "write-only", W: writeOnly{w: io.Discard}},
		{Name: "write-only buffered", W: writeOnly{w: bufio.NewWriterSize(io.Discard, 2048)}},
	}

	attr := func(v struct {
		Name    string
		Boolean bool
		Value   string
	}) g.Node {
		if v.Boolean {
			return g.Attr("hat")
		}
		return g.Attr("hat", v.Value)
	}

	for _, w := range writers {
		for _, v := range values {
			b.Run("construct and render/"+w.Name+"/"+v.Name, func(b *testing.B) {
				for b.Loop() {
					// Nodes in a page are kept by their parent, so keep this one too, or the
					// compiler puts it on the stack and the construction cost disappears.
					node = attr(v)
					_ = node.Render(w.W)
				}
			})

			b.Run("render pre-built/"+w.Name+"/"+v.Name, func(b *testing.B) {
				a := attr(v)

				for b.Loop() {
					_ = a.Render(w.W)
				}
			})
		}
	}
}

func BenchmarkEl(b *testing.B) {
	b.Run("normal elements", func(b *testing.B) {
		for b.Loop() {
			e := g.El("div")
			_ = e.Render(io.Discard)
		}
	})
}

func BenchmarkRaw(b *testing.B) {
	b.Run("raw element", func(b *testing.B) {
		for b.Loop() {
			e := g.Raw("<span>content</span>")
			_ = e.Render(io.Discard)
		}
	})
}

func BenchmarkRawf(b *testing.B) {
	b.Run("formatted raw element", func(b *testing.B) {
		for b.Loop() {
			e := g.Rawf("<span>%s</span>", "content")
			_ = e.Render(io.Discard)
		}
	})
}

func BenchmarkText(b *testing.B) {
	b.Run("simple text element", func(b *testing.B) {
		for b.Loop() {
			e := g.Text("some simple text")
			_ = e.Render(io.Discard)
		}
	})

	b.Run("text element needing escaping", func(b *testing.B) {
		for b.Loop() {
			e := g.Text("It's a sentence & it needs escaping.")
			_ = e.Render(io.Discard)
		}
	})

	b.Run("prose element needing escaping", func(b *testing.B) {
		// English prose contains apostrophes, which on its own is enough to take the
		// escaping path.
		prose := strings.Repeat("It's a paragraph of user-written text, with quotes & apostrophes in it. ", 40)

		for b.Loop() {
			e := g.Text(prose)
			_ = e.Render(io.Discard)
		}
	})
}

func BenchmarkTextf(b *testing.B) {
	b.Run("formatted text element", func(b *testing.B) {
		for b.Loop() {
			e := g.Textf("some %s text", "formatted")
			_ = e.Render(io.Discard)
		}
	})
}
