package gomponents_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/internal/assert"
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
		called := false
		defer func() {
			if err := recover(); err != nil {
				called = true
			}
		}()
		g.Attr("name", "value", "what is this?")
		if !called {
			t.FailNow()
		}
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
	e := g.El("input", g.Attr("required"))
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
		e := g.El("hr")
		assert.Equal(t, "<hr>", e)

		e = g.El("br")
		assert.Equal(t, "<br>", e)

		e = g.El("img")
		assert.Equal(t, "<img>", e)
	})

	t.Run("renders an empty element if only attributes given as children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"></div>`, e)
	})

	t.Run("renders an element, attributes, and element children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"), g.El("br"))
		assert.Equal(t, `<div class="hat"><br></div>`, e)
	})

	t.Run("renders attributes at the correct place regardless of placement in parameter list", func(t *testing.T) {
		e := g.El("div", g.El("br"), g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"><br></div>`, e)
	})

	t.Run("renders outside if node does not implement placer", func(t *testing.T) {
		e := g.El("div", outsider{})
		assert.Equal(t, `<div>outsider</div>`, e)
	})

	t.Run("does not fail on nil node", func(t *testing.T) {
		e := g.El("div", nil, g.El("br"), nil, g.El("br"))
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

func TestGroup(t *testing.T) {
	t.Run("groups multiple nodes into one", func(t *testing.T) {
		children := []g.Node{g.El("br", g.Attr("id", "hat")), g.El("hr")}
		e := g.El("div", g.Attr("class", "foo"), g.El("img"), g.Group(children))
		assert.Equal(t, `<div class="foo"><img><br id="hat"><hr></div>`, e)
	})

	t.Run("panics on direct render", func(t *testing.T) {
		e := g.Group(nil)
		panicked := false
		defer func() {
			if err := recover(); err != nil {
				panicked = true
			}
		}()
		_ = e.Render(nil)
		if !panicked {
			t.FailNow()
		}
	})

	t.Run("panics on direct string", func(t *testing.T) {
		e := g.Group(nil).(fmt.Stringer)
		panicked := false
		defer func() {
			if err := recover(); err != nil {
				panicked = true
			}
		}()
		_ = e.String()
		if !panicked {
			t.FailNow()
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("maps slices to nodes", func(t *testing.T) {
		items := []string{"hat", "partyhat", "turtlehat"}
		lis := g.Map(len(items), func(i int) g.Node {
			return g.El("li", g.Text(items[i]))
		})

		list := g.El("ul", lis...)

		assert.Equal(t, `<ul><li>hat</li><li>partyhat</li><li>turtlehat</li></ul>`, list)
	})
}

func ExampleMap() {
	items := []string{"party hat", "super hat"}
	e := g.El("ul", g.Group(g.Map(len(items), func(i int) g.Node {
		return g.El("li", g.Text(items[i]))
	})))
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>super hat</li></ul>
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

func TestLazy(t *testing.T) {
	t.Run("Lazy is a function type that is also a Node", func(t *testing.T) {
		n := g.El("div",
			g.Lazy(func() g.Node {
				return g.El("span")
			}),
		)
		assert.Equal(t, "<div><span></span></div>", n)
	})

	t.Run("Lazy also works for attribute type nodes", func(t *testing.T) {
		n := g.El("div",
			g.Lazy(func() g.Node {
				return g.Attr("class", "foo")
			}),
		)
		assert.Equal(t, "<input disabled>", n)
	})
}

func ExampleIf_lazy() {
	var message *string
	e := g.El("div",
		g.If(message != nil, g.Lazy(func() g.Node {
			return g.El("span", g.Text(*message))
		})),
	)
	_ = e.Render(os.Stdout)
	// Output: <div></div>
}
