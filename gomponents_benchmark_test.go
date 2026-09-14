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

func BenchmarkAttr(b *testing.B) {
	b.Run("boolean attributes", func(b *testing.B) {
		for b.Loop() {
			a := g.Attr("hat")
			_ = a.Render(io.Discard)
		}
	})

	b.Run("name-value attributes", func(b *testing.B) {
		for b.Loop() {
			a := g.Attr("hat", "party")
			_ = a.Render(io.Discard)
		}
	})

	b.Run("name-value attributes needing escaping", func(b *testing.B) {
		for b.Loop() {
			a := g.Attr("hat", `"party" & fun`)
			_ = a.Render(io.Discard)
		}
	})

	values := []struct {
		Name, Value string
	}{
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

	for _, w := range writers {
		for _, v := range values {
			b.Run("render pre-built/"+w.Name+"/"+v.Name, func(b *testing.B) {
				a := g.Attr("hat", v.Value)

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
