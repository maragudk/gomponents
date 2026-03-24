package gomponents_test

import (
	"html/template"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func FuzzEl(f *testing.F) {
	f.Add("div", "hello")
	f.Add("", "")
	f.Add("br", "child text")
	f.Add("img", "")
	f.Add("<script>", "xss attempt")
	f.Add("div onclick=alert(1)", "")
	f.Add("h1", "Héllo Wörld")
	f.Add("custom-element", "web component")

	f.Fuzz(func(t *testing.T, name, text string) {
		node := g.El(name, g.Text(text))

		var b strings.Builder
		if err := node.Render(&b); err != nil {
			t.Fatal(err)
		}

		got := b.String()

		// Must open with the element name.
		if !strings.HasPrefix(got, "<"+name+">") {
			t.Fatalf("expected prefix %q, got %q", "<"+name+">", got)
		}
	})
}

func FuzzAttr(f *testing.F) {
	f.Add("class", "container")
	f.Add("id", `"quoted"`)
	f.Add("", "")
	f.Add("data-x", "<script>alert(1)</script>")
	f.Add("href", "https://example.com?a=1&b=2")
	f.Add("title", "Héllo & Wörld")

	f.Fuzz(func(t *testing.T, name, value string) {
		node := g.El("div", g.Attr(name, value))

		var b strings.Builder
		if err := node.Render(&b); err != nil {
			t.Fatal(err)
		}

		got := b.String()

		// Attribute value must be escaped.
		escaped := template.HTMLEscapeString(value)
		expected := name + `="` + escaped + `"`
		if !strings.Contains(got, expected) {
			t.Fatalf("expected %q to contain %q", got, expected)
		}
	})
}

func FuzzAttrBool(f *testing.F) {
	f.Add("required")
	f.Add("disabled")
	f.Add("")
	f.Add("data-active")

	f.Fuzz(func(t *testing.T, name string) {
		node := g.El("div", g.Attr(name))

		var b strings.Builder
		if err := node.Render(&b); err != nil {
			t.Fatal(err)
		}

		got := b.String()

		// Boolean attribute must appear as just the name, no ="...".
		if !strings.Contains(got, "<div "+name+">") {
			t.Fatalf("expected %q to contain %q", got, "<div "+name+">")
		}
	})
}

func FuzzText(f *testing.F) {
	f.Add("hello world")
	f.Add("")
	f.Add("<script>alert('xss')</script>")
	f.Add(`"quotes" & 'apostrophes'`)
	f.Add("Ünïcödé \x00\x01\x02")

	f.Fuzz(func(t *testing.T, text string) {
		node := g.Text(text)

		var b strings.Builder
		if err := node.Render(&b); err != nil {
			t.Fatal(err)
		}

		got := b.String()
		expected := template.HTMLEscapeString(text)
		if got != expected {
			t.Fatalf("Text(%q): got %q, expected %q", text, got, expected)
		}
	})
}

func FuzzRaw(f *testing.F) {
	f.Add("hello world")
	f.Add("")
	f.Add("<b>bold</b>")
	f.Add("Ünïcödé \x00\x01\x02")

	f.Fuzz(func(t *testing.T, text string) {
		node := g.Raw(text)

		var b strings.Builder
		if err := node.Render(&b); err != nil {
			t.Fatal(err)
		}

		if b.String() != text {
			t.Fatalf("Raw(%q) modified the input: got %q", text, b.String())
		}
	})
}
