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
