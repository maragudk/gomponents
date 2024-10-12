package gomponents_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/internal/assert"
)

func TestNodeFunc(t *testing.T) {
	t.Run("implements fmt.Stringer", func(t *testing.T) {
		fn := g.NodeFunc(func(w io.Writer) error {
			_, _ = w.Write([]byte("hat"))
			return nil
		})
		if fn.String() != "hat" {
			t.FailNow()
		}
	})
}

func TestAttr(t *testing.T) {
	t.Run("renders just the local name with one argument", func(t *testing.T) {
		a := g.Attr("required")
		assert.Equal(t, " required", a)
	})

	t.Run("renders the name and value when given two arguments", func(t *testing.T) {
		a := g.Attr("id", "hat")
		assert.Equal(t, ` id="hat"`, a)
	})

	t.Run("panics with more than two arguments", func(t *testing.T) {
		defer func() {
			if rec := recover(); rec == nil {
				t.FailNow()
			}
		}()
		g.Attr("name", "value", "what is this?")
	})

	t.Run("implements fmt.Stringer", func(t *testing.T) {
		a := g.Attr("required")
		s := fmt.Sprintf("%v", a)
		if s != " required" {
			t.FailNow()
		}
	})

	t.Run("escapes attribute values", func(t *testing.T) {
		a := g.Attr(`id`, `hat"><script`)
		assert.Equal(t, ` id="hat&#34;&gt;&lt;script"`, a)
	})
}

func BenchmarkAttr(b *testing.B) {
	b.Run("boolean attributes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := g.Attr("hat")
			_ = a.Render(&strings.Builder{})
		}
	})

	b.Run("name-value attributes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := g.Attr("hat", "party")
			_ = a.Render(&strings.Builder{})
		}
	})
}

func ExampleAttr_bool() {
	e := g.VoidEl("input", g.Attr("required"))
	_ = e.Render(os.Stdout)
	// Output: <input required>
}

func ExampleAttr_name_value() {
	e := g.El("div", g.Attr("id", "hat"))
	_ = e.Render(os.Stdout)
	// Output: <div id="hat"></div>
}

type outsider struct{}

func (o outsider) String() string {
	return "outsider"
}

func (o outsider) Render(w io.Writer) error {
	_, _ = w.Write([]byte("outsider"))
	return nil
}

func TestEl(t *testing.T) {
	t.Run("renders an empty element if no children given", func(t *testing.T) {
		e := g.El("div")
		assert.Equal(t, "<div></div>", e)
	})

	t.Run("renders an empty element without closing tag if it's a void kind element", func(t *testing.T) {
		e := g.VoidEl("hr")
		assert.Equal(t, "<hr>", e)

		e = g.VoidEl("br")
		assert.Equal(t, "<br>", e)

		e = g.VoidEl("img")
		assert.Equal(t, "<img>", e)
	})

	t.Run("renders an empty element if only attributes given as children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"></div>`, e)
	})

	t.Run("renders an element, attributes, and element children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"), g.VoidEl("br"))
		assert.Equal(t, `<div class="hat"><br></div>`, e)
	})

	t.Run("renders attributes at the correct place regardless of placement in parameter list", func(t *testing.T) {
		e := g.El("div", g.VoidEl("br"), g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"><br></div>`, e)
	})

	t.Run("renders outside if node does not implement nodeTypeDescriber", func(t *testing.T) {
		e := g.El("div", outsider{})
		assert.Equal(t, `<div>outsider</div>`, e)
	})

	t.Run("does not fail on nil node", func(t *testing.T) {
		e := g.El("div", nil, g.VoidEl("br"), nil, g.VoidEl("br"))
		assert.Equal(t, `<div><br><br></div>`, e)
	})

	t.Run("returns render error on cannot write", func(t *testing.T) {
		e := g.El("div")
		err := e.Render(&erroringWriter{})
		assert.Error(t, err)
	})
}

func BenchmarkEl(b *testing.B) {
	b.Run("normal elements", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e := g.El("div")
			_ = e.Render(&strings.Builder{})
		}
	})
}

func ExampleEl() {
	e := g.El("div", g.El("span"))
	_ = e.Render(os.Stdout)
	// Output: <div><span></span></div>
}

type erroringWriter struct{}

func (w *erroringWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("no thanks")
}

func TestText(t *testing.T) {
	t.Run("renders escaped text", func(t *testing.T) {
		e := g.Text("<div>")
		assert.Equal(t, "&lt;div&gt;", e)
	})
}

func ExampleText() {
	e := g.El("span", g.Text("Party hats > normal hats."))
	_ = e.Render(os.Stdout)
	// Output: <span>Party hats &gt; normal hats.</span>
}

func TestTextf(t *testing.T) {
	t.Run("renders interpolated and escaped text", func(t *testing.T) {
		e := g.Textf("<%v>", "div")
		assert.Equal(t, "&lt;div&gt;", e)
	})
}

func ExampleTextf() {
	e := g.El("span", g.Textf("%v party hats > %v normal hats.", 2, 3))
	_ = e.Render(os.Stdout)
	// Output: <span>2 party hats &gt; 3 normal hats.</span>
}

func TestRaw(t *testing.T) {
	t.Run("renders raw text", func(t *testing.T) {
		e := g.Raw("<div>")
		assert.Equal(t, "<div>", e)
	})
}

func ExampleRaw() {
	e := g.El("span",
		g.Raw(`<button onclick="javascript:alert('Party time!')">Party hats</button> &gt; normal hats.`),
	)
	_ = e.Render(os.Stdout)
	// Output: <span><button onclick="javascript:alert('Party time!')">Party hats</button> &gt; normal hats.</span>
}

func TestRawf(t *testing.T) {
	t.Run("renders interpolated and raw text", func(t *testing.T) {
		e := g.Rawf("<%v>", "div")
		assert.Equal(t, "<div>", e)
	})
}

func ExampleRawf() {
	e := g.El("span",
		g.Rawf(`<button onclick="javascript:alert('%v')">Party hats</button> &gt; normal hats.`, "Party time!"),
	)
	_ = e.Render(os.Stdout)
	// Output: <span><button onclick="javascript:alert('Party time!')">Party hats</button> &gt; normal hats.</span>
}

func TestMap(t *testing.T) {
	t.Run("maps a slice to a group", func(t *testing.T) {
		items := []string{"hat", "partyhat", "turtlehat"}
		lis := g.Map(items, func(i string) g.Node {
			return g.El("li", g.Text(i))
		})

		list := g.El("ul", lis...)

		assert.Equal(t, `<ul><li>hat</li><li>partyhat</li><li>turtlehat</li></ul>`, list)
		if len(lis) != 3 {
			t.FailNow()
		}
		assert.Equal(t, `<li>hat</li>`, lis[0])
		assert.Equal(t, `<li>partyhat</li>`, lis[1])
		assert.Equal(t, `<li>turtlehat</li>`, lis[2])
	})
}

func ExampleMap() {
	items := []string{"party hat", "super hat"}
	e := g.El("ul", g.Map(items, func(i string) g.Node {
		return g.El("li", g.Text(i))
	}))
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>super hat</li></ul>
}

func ExampleMap_index() {
	items := []string{"party hat", "super hat"}
	var index int
	e := g.El("ul", g.Map(items, func(i string) g.Node {
		e := g.El("li", g.Textf("%v: %v", index, i))
		index++
		return e
	}))
	_ = e.Render(os.Stdout)
	// Output: <ul><li>0: party hat</li><li>1: super hat</li></ul>
}

func TestGroup(t *testing.T) {
	t.Run("groups multiple nodes into one", func(t *testing.T) {
		children := []g.Node{g.VoidEl("br", g.Attr("id", "hat")), g.VoidEl("hr")}
		e := g.El("div", g.Attr("class", "foo"), g.VoidEl("img"), g.Group(children))
		assert.Equal(t, `<div class="foo"><img><br id="hat"><hr></div>`, e)
	})

	t.Run("ignores attributes at the first level", func(t *testing.T) {
		children := []g.Node{g.Attr("class", "hat"), g.El("div"), g.El("span")}
		e := g.Group(children)
		assert.Equal(t, "<div></div><span></span>", e)
	})

	t.Run("does not ignore attributes at the second level and below", func(t *testing.T) {
		children := []g.Node{g.El("div", g.Attr("class", "hat"), g.VoidEl("hr", g.Attr("id", "partyhat"))), g.El("span")}
		e := g.Group(children)
		assert.Equal(t, `<div class="hat"><hr id="partyhat"></div><span></span>`, e)
	})

	t.Run("implements fmt.Stringer", func(t *testing.T) {
		children := []g.Node{g.El("div"), g.El("span")}
		e := g.Group(children)
		if e, ok := any(e).(fmt.Stringer); !ok || e.String() != "<div></div><span></span>" {
			t.FailNow()
		}
	})

	t.Run("can be used like a regular slice", func(t *testing.T) {
		e := g.Group{g.El("div"), g.El("span")}
		assert.Equal(t, "<div></div><span></span>", e)
		assert.Equal(t, "<div></div>", e[0])
		assert.Equal(t, "<span></span>", e[1])
	})
}

func ExampleGroup() {
	children := []g.Node{g.El("div"), g.El("span")}
	e := g.Group(children)
	_ = e.Render(os.Stdout)
	// Output: <div></div><span></span>
}

func ExampleGroup_slice() {
	e := g.Group{g.El("div"), g.El("span")}
	_ = e.Render(os.Stdout)
	// Output: <div></div><span></span>
}

func TestIf(t *testing.T) {
	t.Run("returns node if condition is true", func(t *testing.T) {
		n := g.El("div", g.If(true, g.El("span")))
		assert.Equal(t, "<div><span></span></div>", n)
	})

	t.Run("returns nil if condition is false", func(t *testing.T) {
		n := g.El("div", g.If(false, g.El("span")))
		assert.Equal(t, "<div></div>", n)
	})
}

func ExampleIf() {
	showMessage := true

	e := g.El("div",
		g.If(showMessage, g.El("span", g.Text("You lost your hat!"))),
		g.If(!showMessage, g.El("span", g.Text("No messages."))),
	)

	_ = e.Render(os.Stdout)
	// Output: <div><span>You lost your hat!</span></div>
}

func TestIff(t *testing.T) {
	t.Run("returns node if condition is true", func(t *testing.T) {
		n := g.El("div", g.Iff(true, func() g.Node {
			return g.El("span")
		}))
		assert.Equal(t, "<div><span></span></div>", n)
	})

	t.Run("returns nil if condition is false", func(t *testing.T) {
		n := g.El("div", g.Iff(false, func() g.Node {
			return g.El("span")
		}))
		assert.Equal(t, "<div></div>", n)
	})
}

func ExampleIff() {
	type User struct {
		Name string
	}
	var user *User

	e := g.El("div",
		// This would panic using just If
		g.Iff(user != nil, func() g.Node {
			return g.Text(user.Name)
		}),
	)

	_ = e.Render(os.Stdout)
	// Output: <div></div>
}
