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

		// JoinAttrs only joins attributes with non-empty values.
		// Attributes with empty values pass through unmodified.
		var expected string
		switch {
		case value1 != "" && value2 != "":
			expected = renderToString(t, g.El("div", g.Attr(name, value1+" "+value2)))
		default:
			expected = renderToString(t, g.El("div", attr1, attr2))
		}

		if got != expected {
			t.Fatalf("JoinAttrs(%q, %q, %q):\ngot:      %q\nexpected: %q", name, value1, value2, got, expected)
		}
	})
}

func FuzzClasses(f *testing.F) {
	f.Add("hat", "party")
	f.Add("", "")
	f.Add("[&_svg]:size-4", "custom")
	f.Add(`<script>`, `"quoted"`)
	f.Add("a&b", "c<d")

	f.Fuzz(func(t *testing.T, class1, class2 string) {
		node := Classes{class1: true, class2: true}

		var b strings.Builder
		if err := node.Render(&b); err != nil {
			t.Fatal(err)
		}

		got := b.String()

		// Must render as a class attribute.
		if !strings.HasPrefix(got, ` class="`) {
			t.Fatalf("expected class attribute, got %q", got)
		}
		if !strings.HasSuffix(got, `"`) {
			t.Fatalf("expected closing quote, got %q", got)
		}
	})
}
