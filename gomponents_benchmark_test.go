//go:build go1.24

package gomponents_test

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func BenchmarkAttr(b *testing.B) {
	b.Run("boolean attributes", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			a := g.Attr("hat")
			_ = a.Render(&sb)
			sb.Reset()
		}
	})

	b.Run("name-value attributes", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			a := g.Attr("hat", "party")
			_ = a.Render(&sb)
			sb.Reset()
		}
	})
}

func BenchmarkEl(b *testing.B) {
	b.Run("normal elements", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			e := g.El("div")
			_ = e.Render(&sb)
			sb.Reset()
		}
	})
}

func BenchmarkRaw(b *testing.B) {
	b.Run("raw element", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			e := g.Raw("<span>content</span>")
			_ = e.Render(&sb)
			sb.Reset()
		}
	})
}

func BenchmarkRawf(b *testing.B) {
	b.Run("formatted raw element", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			e := g.Rawf("<span>%s</span>", "content")
			_ = e.Render(&sb)
			sb.Reset()
		}
	})
}

func BenchmarkText(b *testing.B) {
	b.Run("simple text element", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			e := g.Text("some simple text")
			_ = e.Render(&sb)
			sb.Reset()
		}
	})
}

func BenchmarkTextf(b *testing.B) {
	b.Run("formatted text element", func(b *testing.B) {
		var sb strings.Builder

		for b.Loop() {
			e := g.Textf("some %s text", "formatted")
			_ = e.Render(&sb)
			sb.Reset()
		}
	})
}
