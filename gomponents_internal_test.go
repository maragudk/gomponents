package gomponents

import (
	"html/template"
	"strings"
	"testing"
)

func TestRaw(t *testing.T) {
	t.Run("r.String() == string(r)", func(t *testing.T) {
		r := raw("<p>raw</p>")
		if r.String() != string(r) {
			t.Fail()
		}
	})
}

func TestEscapeString(t *testing.T) {
	// escapeString must agree with template.HTMLEscapeString on every input. It no longer
	// calls it — it decides whether to escape and then escapes on its own — so this is what
	// holds the two together. A false negative would pass markup through.
	same := func(t *testing.T, s string) {
		t.Helper()
		if got, want := escapeString(s), template.HTMLEscapeString(s); got != want {
			t.Fatalf("escapeString(%q) = %q, template.HTMLEscapeString = %q", s, got, want)
		}
	}

	t.Run("every single byte", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			same(t, string([]byte{byte(i)}))
		}
	})

	t.Run("every pair of bytes", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				same(t, string([]byte{byte(i), byte(j)}))
			}
		}
	})

	t.Run("longer than the eight bytes that switch strings.IndexAny strategy", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			same(t, "0123456789"+string([]byte{byte(i)})+"0123456789")
		}
	})

	t.Run("either side of the forty-eight bytes that switch strings.Replacer strategy", func(t *testing.T) {
		// Below len(toReplace)*countCutOff, which is 48 for six characters,
		// byteStringReplacer sizes its output by walking the string once; from there on it
		// counts each of the six separately. Both paths have to agree with the reference.
		for _, n := range []int{47, 48, 49} {
			for _, c := range []string{"\x00", `"`, "'", "&", "<", ">"} {
				same(t, c+strings.Repeat("a", n-1))
				same(t, strings.Repeat("a", n-1)+c)
				same(t, strings.Repeat("a", n/2)+c+strings.Repeat("a", n-n/2-1))
			}
		}
	})

	t.Run("long strings on the counting path", func(t *testing.T) {
		same(t, "&"+strings.Repeat("a", 4095))
		same(t, strings.Repeat("a", 4095)+"&")
		same(t, strings.Repeat(`a"b'c&d<e>f`, 400))
		same(t, strings.Repeat("\x00", 4096))
	})

	t.Run("each escaped character on its own", func(t *testing.T) {
		for _, s := range []string{"\x00", `"`, "'", "&", "<", ">"} {
			same(t, s)
		}
	})

	t.Run("strings worth naming", func(t *testing.T) {
		for _, s := range []string{
			"", " ", "hat", "party hat",
			"<script>", "a&b", `"quoted"`, "'single'", ">", "<", "&", "\x00",
			"&amp;", "&lt;script&gt;",
			"héj", "日本語", "🎉", "héj & 日本語 <b>",
			"a" + string([]byte{0}) + "b",
		} {
			same(t, s)
		}
	})
}

func FuzzEscapeString(f *testing.F) {
	for _, s := range []string{"", "hat", "<script>", "a&b", `"x"`, "日本語", "\x00"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if got, want := escapeString(s), template.HTMLEscapeString(s); got != want {
			t.Fatalf("escapeString(%q) = %q, template.HTMLEscapeString = %q", s, got, want)
		}
	})
}
