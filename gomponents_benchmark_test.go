//go:build go1.24

package gomponents_test

import (
	"bufio"
	"io"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

// value returns n bytes of plain text with the given number of bytes, spread evenly,
// swapped for ones that need escaping, cycling through the three most common.
func value(n, escapes int) string {
	b := []byte(strings.Repeat("party hat ", n/10+1)[:n])
	for i := 0; i < escapes; i++ {
		b[(i+1)*n/(escapes+1)] = `"&'`[i%3]
	}
	return string(b)
}

// writers to render to, each constructed anew by New so that a buffer's fill level cannot
// carry over from one sub-benchmark to the next. The buffered writers are the size of the
// [bufio.Writer] that net/http puts in front of a handler's response writer, its
// bufferBeforeChunkingSize.
func writers() []struct {
	Name string
	New  func() io.Writer
} {
	return []struct {
		Name string
		New  func() io.Writer
	}{
		{Name: "discarded", New: func() io.Writer { return io.Discard }},
		{Name: "buffered", New: func() io.Writer { return bufio.NewWriterSize(io.Discard, 2048) }},
		{Name: "write-only", New: func() io.Writer { return writeOnly{w: io.Discard} }},
		{Name: "write-only buffered", New: func() io.Writer { return writeOnly{w: bufio.NewWriterSize(io.Discard, 2048)} }},
	}
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

	for _, w := range writers() {
		for _, v := range values {
			b.Run("construct and render/"+w.Name+"/"+v.Name, func(b *testing.B) {
				// Nodes in a page are kept by their parent, so keep this one too, or the
				// compiler puts it on the stack and the construction cost disappears.
				var node g.Node
				w := w.New()

				for b.Loop() {
					node = attr(v)
					_ = node.Render(w)
				}
			})

			b.Run("render pre-built/"+w.Name+"/"+v.Name, func(b *testing.B) {
				a := attr(v)
				w := w.New()

				for b.Loop() {
					_ = a.Render(w)
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
	// Short and long text at each level of escaping.
	values := []struct {
		Name, Value string
	}{
		{Name: "short text, no escaping", Value: value(16, 0)},
		{Name: "long text, no escaping", Value: value(256, 0)},
		{Name: "short text, little escaping", Value: value(16, 1)},
		{Name: "long text, little escaping", Value: value(256, 16)},
		{Name: "short text, much escaping", Value: value(16, 4)},
		{Name: "long text, much escaping", Value: value(256, 64)},
	}

	for _, w := range writers() {
		for _, v := range values {
			b.Run("construct and render/"+w.Name+"/"+v.Name, func(b *testing.B) {
				// Nodes in a page are kept by their parent, so keep this one too, or the
				// compiler puts it on the stack and the construction cost disappears.
				var node g.Node
				w := w.New()

				for b.Loop() {
					node = g.Text(v.Value)
					_ = node.Render(w)
				}
			})

			b.Run("render pre-built/"+w.Name+"/"+v.Name, func(b *testing.B) {
				t := g.Text(v.Value)
				w := w.New()

				for b.Loop() {
					_ = t.Render(w)
				}
			})
		}
	}
}

func BenchmarkTextf(b *testing.B) {
	b.Run("formatted text element", func(b *testing.B) {
		for b.Loop() {
			e := g.Textf("some %s text", "formatted")
			_ = e.Render(io.Discard)
		}
	})
}
