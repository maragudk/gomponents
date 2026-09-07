//go:build go1.24

package components_test

import (
	"io"
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
