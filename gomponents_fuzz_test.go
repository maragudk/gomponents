package gomponents_test

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

// onlyWriter wraps an io.Writer to hide the io.StringWriter interface,
// forcing the non-StringWriter code path.
type onlyWriter struct {
	buf bytes.Buffer
}

func (w *onlyWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func FuzzEl(f *testing.F) {
	f.Add("div", "hello")
	f.Add("", "")
	f.Add("br", "ignored for void")
	f.Add("img", "")
	f.Add("<script>", "xss attempt")
	f.Add("div onclick=alert(1)", "")
	f.Add("h1", "Héllo Wörld")
	f.Add("custom-element", "web component")

	f.Fuzz(func(t *testing.T, name, text string) {
		node := g.El(name, g.Text(text))

		// Render to a StringWriter (strings.Builder).
		var sw strings.Builder
		err1 := node.Render(&sw)

		// Render to a plain io.Writer (no StringWriter interface).
		var ow onlyWriter
		err2 := node.Render(&ow)

		// Both code paths must agree on error/success.
		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("error mismatch: StringWriter err = %v, Writer err = %v", err1, err2)
		}

		// Both code paths must produce identical output.
		if sw.String() != ow.buf.String() {
			t.Fatalf("output mismatch for El(%q, Text(%q)):\nStringWriter: %q\nWriter:       %q",
				name, text, sw.String(), ow.buf.String())
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
		node := g.El("div", g.Attr(name, value), g.Text("x"))

		var sw strings.Builder
		err1 := node.Render(&sw)

		var ow onlyWriter
		err2 := node.Render(&ow)

		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("error mismatch: StringWriter err = %v, Writer err = %v", err1, err2)
		}

		if sw.String() != ow.buf.String() {
			t.Fatalf("output mismatch for Attr(%q, %q):\nStringWriter: %q\nWriter:       %q",
				name, value, sw.String(), ow.buf.String())
		}
	})
}

func FuzzText(f *testing.F) {
	f.Add("hello world")
	f.Add("")
	f.Add("<script>alert('xss')</script>")
	f.Add(`"quotes" & 'apostrophes'`)
	f.Add("Ünïcödé \x00\x01\x02")
	f.Add(strings.Repeat("a", 10000))

	f.Fuzz(func(t *testing.T, text string) {
		node := g.Text(text)

		var sw strings.Builder
		err1 := node.Render(&sw)

		var ow onlyWriter
		err2 := node.Render(&ow)

		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("error mismatch: StringWriter err = %v, Writer err = %v", err1, err2)
		}

		if sw.String() != ow.buf.String() {
			t.Fatalf("output mismatch for Text(%q):\nStringWriter: %q\nWriter:       %q",
				text, sw.String(), ow.buf.String())
		}

		// Text must never contain unescaped HTML special characters from the input.
		rendered := sw.String()
		if strings.ContainsAny(text, "<>&\"'") && rendered == text {
			t.Fatalf("Text(%q) was not escaped: got %q", text, rendered)
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

		var sw strings.Builder
		err1 := node.Render(&sw)

		var ow onlyWriter
		err2 := node.Render(&ow)

		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("error mismatch: StringWriter err = %v, Writer err = %v", err1, err2)
		}

		if sw.String() != ow.buf.String() {
			t.Fatalf("output mismatch for Raw(%q):\nStringWriter: %q\nWriter:       %q",
				text, sw.String(), ow.buf.String())
		}

		// Raw must render the exact input, unmodified.
		if sw.String() != text {
			t.Fatalf("Raw(%q) modified the input: got %q", text, sw.String())
		}
	})
}
