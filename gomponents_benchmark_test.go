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

// value returns n bytes of plain text with the given number of bytes, spread evenly,
// swapped for ones that need escaping, cycling through the three most common.
func value(n, escapes int) string {
	b := []byte(strings.Repeat("party hat ", n/10+1)[:n])
	for i := 0; i < escapes; i++ {
		b[(i+1)*n/(escapes+1)] = `"&'`[i%3]
	}
	return string(b)
}

func BenchmarkAttr(b *testing.B) {
	// A boolean attribute has no value; the rest are name-value attributes with short and
	// long values at each level of escaping.
	values := []struct {
		Name    string
		Boolean bool
		Value   string
	}{
		{Name: "boolean", Boolean: true},
		{Name: "short value, no escaping", Value: value(16, 0)},
		{Name: "long value, no escaping", Value: value(256, 0)},
		{Name: "short value, little escaping", Value: value(16, 1)},
		{Name: "long value, little escaping", Value: value(256, 16)},
		{Name: "short value, much escaping", Value: value(16, 4)},
		{Name: "long value, much escaping", Value: value(256, 64)},
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
				// Nodes in a page are kept by their parent, so keep this one too, or the
				// compiler puts it on the stack and the construction cost disappears.
				var node g.Node

				for b.Loop() {
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
