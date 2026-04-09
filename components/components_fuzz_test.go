package components_test

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
)

func renderToString(t *testing.T, n g.Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func FuzzJoinAttrs(f *testing.F) {
	f.Add("class", "party", "hat")
	f.Add("class", "", "")
	f.Add("class", "[&_svg]:size-4", "custom")
	f.Add("data-test", `<script>"test"</script>`, "more")
	f.Add("class", `"quotes"`, `'apostrophes'`)
	f.Add("class", "a&b", "c&amp;d")
	f.Add("class", "a&lt;b", "c<d")
	f.Add("style", "color: red", "font-size: 12px")

	f.Fuzz(func(t *testing.T, name, value1, value2 string) {
		// Create two attributes and join them, rendered inside an element.
		attr1 := g.Attr(name, value1)
		attr2 := g.Attr(name, value2)
		got := renderToString(t, g.El("div", JoinAttrs(name, attr1, attr2)))

		// JoinAttrs should always produce a single merged attribute.
		// See https://github.com/maragudk/gomponents/issues/302
		// JoinAttrs treats whitespace-only values as empty.
		v1 := strings.TrimSpace(value1)
		v2 := strings.TrimSpace(value2)
		var expected string
		switch {
		case v1 != "" && v2 != "":
			expected = renderToString(t, g.El("div", g.Attr(name, value1+" "+value2)))
		case v1 != "":
			expected = renderToString(t, g.El("div", g.Attr(name, value1)))
		case v2 != "":
			expected = renderToString(t, g.El("div", g.Attr(name, value2)))
		default:
			expected = renderToString(t, g.El("div", g.Attr(name)))
		}

		if got != expected {
			t.Fatalf("JoinAttrs(%q, %q, %q):\ngot:      %q\nexpected: %q", name, value1, value2, got, expected)
		}
	})
}

