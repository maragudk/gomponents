package gomponents_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
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
		assert.Equal(t, "<div />", e)
	})

	t.Run("renders an empty element if only attributes given as children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat" />`, e)
	})

	t.Run("renders an element, attributes, and element children", func(t *testing.T) {
		e := g.El("div", g.Attr("class", "hat"), g.El("span"))
		assert.Equal(t, `<div class="hat"><span /></div>`, e)
	})

	t.Run("renders attributes at the correct place regardless of placement in parameter list", func(t *testing.T) {
		e := g.El("div", g.El("span"), g.Attr("class", "hat"))
		assert.Equal(t, `<div class="hat"><span /></div>`, e)
	})

	t.Run("renders outside if node does not implement placer", func(t *testing.T) {
		e := g.El("div", outsider{})
		assert.Equal(t, `<div>outsider</div>`, e)
	})

	t.Run("does not fail on nil node", func(t *testing.T) {
		e := g.El("div", nil, g.El("span"), nil, g.El("span"))
		assert.Equal(t, `<div><span /><span /></div>`, e)
	})

	t.Run("returns render error on cannot write", func(t *testing.T) {
		e := g.El("div")
		err := e.Render(&erroringWriter{})
		assert.Error(t, err)
	})
}

type erroringWriter struct{}

func (w *erroringWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("no thanks")
}

func TestText(t *testing.T) {
	t.Run("renders escaped text", func(t *testing.T) {
		e := g.Text("<div />")
		assert.Equal(t, "&lt;div /&gt;", e)
	})
}

func TestTextf(t *testing.T) {
	t.Run("renders interpolated and escaped text", func(t *testing.T) {
		e := g.Textf("<%v />", "div")
		assert.Equal(t, "&lt;div /&gt;", e)
	})
}

func TestRaw(t *testing.T) {
	t.Run("renders raw text", func(t *testing.T) {
		e := g.Raw("<div />")
		assert.Equal(t, "<div />", e)
	})
}

func TestGroup(t *testing.T) {
	t.Run("groups multiple nodes into one", func(t *testing.T) {
		children := []g.Node{g.El("div", g.Attr("id", "hat")), g.El("div")}
		e := g.El("div", g.Attr("class", "foo"), g.El("div"), g.Group(children))
		assert.Equal(t, `<div class="foo"><div /><div id="hat" /><div /></div>`, e)
	})

	t.Run("panics on direct render", func(t *testing.T) {
		e := g.Group(nil)
		panicced := false
		defer func() {
			if err := recover(); err != nil {
				panicced = true
			}
		}()
		_ = e.Render(nil)
		if !panicced {
			t.FailNow()
		}
	})

	t.Run("panics on direct string", func(t *testing.T) {
		e := g.Group(nil)
		panicced := false
		defer func() {
			if err := recover(); err != nil {
				panicced = true
			}
		}()
		_ = fmt.Sprintf("%v", e)
		if !panicced {
			t.FailNow()
		}
	})
}
