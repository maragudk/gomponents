package components_test

import (
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	"maragu.dev/gomponents/internal/assert"
)

func TestHTML5(t *testing.T) {
	t.Run("returns an html5 document template", func(t *testing.T) {
		e := HTML5(HTML5Props{
			Title:       "Hat",
			Description: "Love hats.",
			Language:    "en",
			Head:        []g.Node{Link(Rel("stylesheet"), Href("/hat.css"))},
			Body:        []g.Node{Div()},
		})

		assert.Equal(t, `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Hat</title><meta name="description" content="Love hats."><link rel="stylesheet" href="/hat.css"></head><body><div></div></body></html>`, e)
	})

	t.Run("returns no language, description, and extra head/body elements if empty", func(t *testing.T) {
		e := HTML5(HTML5Props{
			Title: "Hat",
		})

		assert.Equal(t, `<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Hat</title></head><body></body></html>`, e)
	})

	t.Run("returns an html5 document template with additional HTML attributes", func(t *testing.T) {
		e := HTML5(HTML5Props{
			Title:       "Hat",
			Description: "Love hats.",
			Language:    "en",
			Head:        []g.Node{Link(Rel("stylesheet"), Href("/hat.css"))},
			Body:        []g.Node{Div()},
			HTMLAttrs:   []g.Node{Class("h-full"), ID("htmlid")},
		})

		assert.Equal(t, `<!doctype html><html lang="en" class="h-full" id="htmlid"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Hat</title><meta name="description" content="Love hats."><link rel="stylesheet" href="/hat.css"></head><body><div></div></body></html>`, e)
	})

	t.Run("accepts g.Group literal syntax", func(t *testing.T) {
		e := HTML5(HTML5Props{
			Title:     "Hat",
			Head:      g.Group{Link(Rel("stylesheet"), Href("/hat.css"))},
			Body:      g.Group{Div()},
			HTMLAttrs: g.Group{Class("h-full")},
		})

		assert.Equal(t, `<!doctype html><html class="h-full"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Hat</title><link rel="stylesheet" href="/hat.css"></head><body><div></div></body></html>`, e)
	})
}

func TestClasses(t *testing.T) {
	t.Run("given a map, returns sorted keys from the map with value true", func(t *testing.T) {
		assert.Equal(t, ` class="boheme-hat hat partyhat"`, Classes{
			"boheme-hat": true,
			"hat":        true,
			"partyhat":   true,
			"turtlehat":  false,
		})
	})

	t.Run("renders as attribute in an element", func(t *testing.T) {
		e := g.El("div", Classes{"hat": true})
		assert.Equal(t, `<div class="hat"></div>`, e)
	})

	t.Run("also works with fmt", func(t *testing.T) {
		a := Classes{"hat": true}
		if a.String() != ` class="hat"` {
			t.FailNow()
		}
	})
}

func ExampleClasses() {
	e := g.El("div", Classes{"party-hat": true, "boring-hat": false})
	_ = e.Render(os.Stdout)
	// Output: <div class="party-hat"></div>
}

func hat(children ...g.Node) g.Node {
	return Div(JoinAttrs("class", g.Group(children), Class("hat")))
}

func partyHat(children ...g.Node) g.Node {
	return hat(ID("party-hat"), Class("party"), g.Group(children))
}

type brokenNode struct {
	first bool
}

func (b *brokenNode) Render(io.Writer) error {
	if !b.first {
		return nil
	}
	b.first = false
	return errors.New("oh no")
}

func (b *brokenNode) Type() g.NodeType {
	return g.AttributeType
}

// truncatedAttrNode renders as an attribute whose value is cut off after the opening quote,
// which is the one shape where the rendered text ends with the same quote it starts the value with.
type truncatedAttrNode struct{ name string }

func (t truncatedAttrNode) Render(w io.Writer) error {
	_, err := io.WriteString(w, " "+t.name+`="`)
	return err
}

func (truncatedAttrNode) Type() g.NodeType { return g.AttributeType }

// recorder is a [g.Node] that records whether Render was called on it.
type recorder struct{ rendered bool }

func (r *recorder) Render(out io.Writer) error {
	r.rendered = true
	_, err := io.WriteString(out, "!")
	return err
}

// elementNode reports Type as [g.ElementType].
type elementNode struct{ *recorder }

func (elementNode) Type() g.NodeType { return g.ElementType }

// defaultTypeNode has no Type method, so it defaults to [g.ElementType] per the [g.NodeType] contract.
type defaultTypeNode struct{ *recorder }

func TestJoinAttrs(t *testing.T) {
	t.Run("joins classes", func(t *testing.T) {
		n := Div(JoinAttrs("class", Class("party"), ID("hey"), Class("hat")))
		assert.Equal(t, `<div class="party hat" id="hey"></div>`, n)
	})

	t.Run("joins classes in groups", func(t *testing.T) {
		n := partyHat(Span(ID("party-hat-text"), Class("solid"), Class("gold"), g.Text("Yo.")))
		assert.Equal(t, `<div id="party-hat" class="party hat"><span id="party-hat-text" class="solid" class="gold">Yo.</span></div>`, n)
	})

	t.Run("does nothing if attribute not found", func(t *testing.T) {
		n := Div(JoinAttrs("style", Class("party"), ID("hey"), Class("hat")))
		assert.Equal(t, `<div class="party" id="hey" class="hat"></div>`, n)
	})

	t.Run("treats an attribute truncated after the opening quote as having no value", func(t *testing.T) {
		n := Div(JoinAttrs("class", truncatedAttrNode{name: "class"}))
		assert.Equal(t, `<div class></div>`, n)
	})

	t.Run("lets a real value win over a truncated one", func(t *testing.T) {
		n := Div(JoinAttrs("class", truncatedAttrNode{name: "class"}, Class("hat")))
		assert.Equal(t, `<div class="hat"></div>`, n)
	})

	t.Run("ignores nodes that can't render", func(t *testing.T) {
		n := Div(JoinAttrs("class", Class("party"), ID("hey"), &brokenNode{first: true}, Class("hat")))
		assert.Equal(t, `<div class="party hat" id="hey"></div>`, n)
	})

	t.Run("does not eagerly render nodes classified as elements while inspecting attributes", func(t *testing.T) {
		kinds := []struct {
			name    string
			newNode func() (g.Node, *recorder)
		}{
			{"element-typed node", func() (g.Node, *recorder) { r := &recorder{}; return elementNode{r}, r }},
			{"default-typed node", func() (g.Node, *recorder) { r := &recorder{}; return defaultTypeNode{r}, r }},
		}

		scenarios := []struct {
			name     string
			attrName string
			children func(n g.Node) []g.Node
			expected string
		}{
			{"before the matching attribute", "class", func(n g.Node) []g.Node { return []g.Node{n, Class("party"), Class("hat")} }, `<div class="party hat">!</div>`},
			{"nested one group deep", "class", func(n g.Node) []g.Node { return []g.Node{Class("party"), g.Group{n}} }, `<div class="party">!</div>`},
			{"nested two groups deep", "class", func(n g.Node) []g.Node { return []g.Node{Class("party"), g.Group{g.Group{n}}} }, `<div class="party">!</div>`},
			{"when no attribute matches", "required", func(n g.Node) []g.Node { return []g.Node{n} }, `<div>!</div>`},
		}

		for _, kind := range kinds {
			for _, scenario := range scenarios {
				t.Run(kind.name+" "+scenario.name, func(t *testing.T) {
					node, rec := kind.newNode()

					result := JoinAttrs(scenario.attrName, scenario.children(node)...)
					if rec.rendered {
						t.Fatal("node was rendered while JoinAttrs was still inspecting attributes")
					}

					assert.Equal(t, scenario.expected, Div(result))
					if !rec.rendered {
						t.Fatal("node was never rendered, so this test proves nothing")
					}
				})
			}
		}
	})

	t.Run("discards empty-valued matching attributes", func(t *testing.T) {
		n := Div(JoinAttrs("class", Class("party"), g.Attr("class", "")))
		assert.Equal(t, `<div class="party"></div>`, n)
	})

	t.Run("discards whitespace-only matching attributes", func(t *testing.T) {
		n := Div(JoinAttrs("class", Class("party"), g.Attr("class", "  ")))
		assert.Equal(t, `<div class="party"></div>`, n)
	})

	t.Run("discards empty-valued matching attributes in groups", func(t *testing.T) {
		n := Div(JoinAttrs("class", g.Group{Class("party"), g.Attr("class", "")}))
		assert.Equal(t, `<div class="party"></div>`, n)
	})

	t.Run("deduplicates boolean attributes", func(t *testing.T) {
		n := Div(JoinAttrs("required", g.Attr("required"), ID("hey"), g.Attr("required")))
		assert.Equal(t, `<div required id="hey"></div>`, n)
	})

	t.Run("deduplicates boolean attributes in groups", func(t *testing.T) {
		n := Div(JoinAttrs("required", g.Group{g.Attr("required"), g.Attr("required")}))
		assert.Equal(t, `<div required></div>`, n)
	})

	t.Run("keeps single boolean attribute", func(t *testing.T) {
		n := Div(JoinAttrs("required", g.Attr("required"), ID("hey")))
		assert.Equal(t, `<div required id="hey"></div>`, n)
	})

	t.Run("valued attribute takes precedence over boolean", func(t *testing.T) {
		n := Div(JoinAttrs("hidden", g.Attr("hidden"), g.Attr("hidden", "until-found")))
		assert.Equal(t, `<div hidden="until-found"></div>`, n)
	})

	t.Run("does not double-escape ampersands", func(t *testing.T) {
		n := Div(JoinAttrs("class", Class("[&_svg]:size-4"), Class("custom")))
		assert.Equal(t, `<div class="[&amp;_svg]:size-4 custom"></div>`, n)
	})

	t.Run("does not double-escape other HTML entities", func(t *testing.T) {
		n := Div(JoinAttrs("data-test", g.Attr("data-test", `<script>"test"</script>`), g.Attr("data-test", "more")))
		assert.Equal(t, `<div data-test="&lt;script&gt;&#34;test&#34;&lt;/script&gt; more"></div>`, n)
	})

	t.Run("joins classes in nested groups", func(t *testing.T) {
		n := Div(JoinAttrs("class", g.Group{Class("party"), g.Group{Class("gold")}}, Class("hat")))
		assert.Equal(t, `<div class="party gold hat"></div>`, n)
	})

	t.Run("joins classes in deeply nested groups", func(t *testing.T) {
		n := Div(JoinAttrs("class", g.Group{g.Group{g.Group{Class("a")}, Class("b")}, Class("c")}, Class("d")))
		assert.Equal(t, `<div class="a b c d"></div>`, n)
	})

	t.Run("joins classes passed through nested components", func(t *testing.T) {
		n := partyHat(Class("gold"), g.Text("Yo."))
		assert.Equal(t, `<div id="party-hat" class="party gold hat">Yo.</div>`, n)
	})

	t.Run("keeps element order when joining classes in nested groups", func(t *testing.T) {
		n := Div(JoinAttrs("class", g.Group{g.Text("a"), g.Group{Class("x"), g.Text("b")}, Class("y")}))
		assert.Equal(t, `<div class="x y">ab</div>`, n)
	})
}

// failingWriter fails every write.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("no thanks")
}

func TestStatic(t *testing.T) {
	t.Run("calls f and renders the node once and reuses the cached HTML on later renders", func(t *testing.T) {
		var slot string
		var calls int32
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return P(Class("hat"), g.Text("Party hat"))
		}

		// A new Static call each time, like a component called per request.
		for i := 0; i < 3; i++ {
			assert.Equal(t, `<p class="hat">Party hat</p>`, Static(&slot, f))
		}

		if calls != 1 {
			t.Fatalf("expected 1 call, got %v", calls)
		}
		if slot != `<p class="hat">Party hat</p>` {
			t.Fatalf("expected the slot to hold the HTML, got %q", slot)
		}
	})

	t.Run("returns a render error, leaves the slot empty, and calls f again on the next render", func(t *testing.T) {
		var slot string
		var calls int32
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return g.NodeFunc(func(w io.Writer) error {
				return errors.New("oh no")
			})
		}

		for i := 0; i < 2; i++ {
			err := Static(&slot, f).Render(io.Discard)
			assert.Error(t, err)
			if slot != "" {
				t.Fatalf("expected an empty slot after a render error, got %q", slot)
			}
		}

		if calls != 2 {
			t.Fatalf("expected 2 calls, got %v", calls)
		}
	})

	t.Run("returns a write error but keeps the cache and does not call f again", func(t *testing.T) {
		var slot string
		var calls int32
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return P(g.Text("hat"))
		}

		err := Static(&slot, f).Render(failingWriter{})
		assert.Error(t, err)
		if slot != "<p>hat</p>" {
			t.Fatalf("expected the slot to be filled despite the write error, got %q", slot)
		}

		assert.Equal(t, "<p>hat</p>", Static(&slot, f))
		if calls != 1 {
			t.Fatalf("expected 1 call, got %v", calls)
		}
	})

	t.Run("calls f every time when the node renders nothing", func(t *testing.T) {
		var slot string
		var calls int32
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return g.Group{}
		}

		for i := 0; i < 3; i++ {
			assert.Equal(t, "", Static(&slot, f))
		}

		if calls != 3 {
			t.Fatalf("expected 3 calls, got %v", calls)
		}
	})

	t.Run("renders nothing when f returns nil", func(t *testing.T) {
		var slot string

		assert.Equal(t, "<div></div>", Div(Static(&slot, func() g.Node { return g.If(false, Span()) })))
		if slot != "" {
			t.Fatalf("expected an empty slot, got %q", slot)
		}
	})

	t.Run("renders whichever node came first when two call sites share a slot", func(t *testing.T) {
		var slot string

		assert.Equal(t, "<p>first</p>", Static(&slot, func() g.Node { return P(g.Text("first")) }))
		assert.Equal(t, "<p>first</p>", Static(&slot, func() g.Node { return P(g.Text("second")) }))
	})

	t.Run("calls f again after the slot is reset", func(t *testing.T) {
		var slot string
		var calls int32
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return P(g.Text("hat"))
		}

		assert.Equal(t, "<p>hat</p>", Static(&slot, f))
		slot = ""
		assert.Equal(t, "<p>hat</p>", Static(&slot, f))

		if calls != 2 {
			t.Fatalf("expected 2 calls, got %v", calls)
		}
	})

	t.Run("renders as an element, not an attribute", func(t *testing.T) {
		var slot string

		assert.Equal(t, `<div> class="hat"</div>`, Div(Static(&slot, func() g.Node { return Class("hat") })))
	})

	t.Run("fills both slots when one Static is nested in another", func(t *testing.T) {
		var outer, inner string

		assert.Equal(t, "<div><span>hat</span></div>", Static(&outer, func() g.Node {
			return Div(Static(&inner, func() g.Node { return Span(g.Text("hat")) }))
		}))
		if outer != "<div><span>hat</span></div>" {
			t.Fatalf("expected the outer slot to be filled, got %q", outer)
		}
		if inner != "<span>hat</span>" {
			t.Fatalf("expected the inner slot to be filled, got %q", inner)
		}
	})

	t.Run("leaves the slot empty when the node panics", func(t *testing.T) {
		var slot string
		f := func() g.Node {
			return g.NodeFunc(func(w io.Writer) error {
				panic("oh no")
			})
		}

		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expected a panic")
				}
			}()
			_ = Static(&slot, f).Render(io.Discard)
		}()

		if slot != "" {
			t.Fatalf("expected an empty slot, got %q", slot)
		}
		assert.Equal(t, "<p>hat</p>", Static(&slot, func() g.Node { return P(g.Text("hat")) }))
	})

	t.Run("fills the slot and gives every goroutine the full output when many render it concurrently", func(t *testing.T) {
		const goroutines = 64

		var slot string
		var calls int32
		start := make(chan struct{})
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return P(g.Text("hat"))
		}

		outputs := make([]string, goroutines)
		errs := make([]error, goroutines)
		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start
				var b strings.Builder
				errs[i] = Static(&slot, f).Render(&b)
				outputs[i] = b.String()
			}(i)
		}
		close(start)
		wg.Wait()

		for i := 0; i < goroutines; i++ {
			if errs[i] != nil {
				t.Fatalf("goroutine %v got error %v", i, errs[i])
			}
			if outputs[i] != "<p>hat</p>" {
				t.Fatalf("goroutine %v got %q", i, outputs[i])
			}
		}
		if slot != "<p>hat</p>" {
			t.Fatalf("expected the slot to hold the HTML, got %q", slot)
		}
		if calls < 1 {
			t.Fatal("expected at least 1 call")
		}
	})

	t.Run("lets readers of a filled slot proceed while another slot's first render is in progress", func(t *testing.T) {
		const readers = 32

		var slotA, slotB string
		var callsA int32
		fA := func() g.Node {
			atomic.AddInt32(&callsA, 1)
			return P(g.Text("a"))
		}
		assert.Equal(t, "<p>a</p>", Static(&slotA, fA))

		started := make(chan struct{})
		release := make(chan struct{})
		fB := func() g.Node {
			return g.NodeFunc(func(w io.Writer) error {
				close(started)
				<-release
				_, err := io.WriteString(w, "<p>b</p>")
				return err
			})
		}

		var writer sync.WaitGroup
		writer.Add(1)
		var outputB string
		var errB error
		go func() {
			defer writer.Done()
			var b strings.Builder
			errB = Static(&slotB, fB).Render(&b)
			outputB = b.String()
		}()
		// Now the first render of slot B is in progress, inside its node, until release is closed.
		<-started

		outputs := make([]string, readers)
		errs := make([]error, readers)
		var wg sync.WaitGroup
		for i := 0; i < readers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				var b strings.Builder
				errs[i] = Static(&slotA, fA).Render(&b)
				outputs[i] = b.String()
			}(i)
		}
		// The readers finish while slot B is still rendering; this would hang if they had to wait for it.
		wg.Wait()
		close(release)
		writer.Wait()

		if errB != nil {
			t.Fatal("writer got error:", errB)
		}
		if outputB != "<p>b</p>" {
			t.Fatalf("writer got %q", outputB)
		}
		for i := 0; i < readers; i++ {
			if errs[i] != nil {
				t.Fatalf("reader %v got error %v", i, errs[i])
			}
			if outputs[i] != "<p>a</p>" {
				t.Fatalf("reader %v got %q", i, outputs[i])
			}
		}
		if callsA != 1 {
			t.Fatalf("expected 1 call for a, got %v", callsA)
		}
	})

	t.Run("keeps the first result when concurrent first renders of a slot all try to store it", func(t *testing.T) {
		// Every goroutine sees the empty slot and calls f, since the node blocks until all of them
		// have. They then all try to store their result, and only the first one does.
		const goroutines = 32

		var slot string
		var calls int32
		release := make(chan struct{})
		f := func() g.Node {
			atomic.AddInt32(&calls, 1)
			return g.NodeFunc(func(w io.Writer) error {
				<-release
				_, err := io.WriteString(w, "<p>hat</p>")
				return err
			})
		}

		outputs := make([]string, goroutines)
		errs := make([]error, goroutines)
		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				var b strings.Builder
				errs[i] = Static(&slot, f).Render(&b)
				outputs[i] = b.String()
			}(i)
		}
		for atomic.LoadInt32(&calls) < goroutines {
			runtime.Gosched()
		}
		close(release)
		wg.Wait()

		for i := 0; i < goroutines; i++ {
			if errs[i] != nil {
				t.Fatalf("goroutine %v got error %v", i, errs[i])
			}
			if outputs[i] != "<p>hat</p>" {
				t.Fatalf("goroutine %v got %q", i, outputs[i])
			}
		}
		if slot != "<p>hat</p>" {
			t.Fatalf("expected the slot to hold the HTML, got %q", slot)
		}
		if calls != goroutines {
			t.Fatalf("expected %v calls, got %v", goroutines, calls)
		}
	})
}
