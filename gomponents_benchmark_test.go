//go:build go1.24

package gomponents_test

import (
	"io"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

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
