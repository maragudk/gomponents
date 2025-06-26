package components_test

import (
	"errors"
	"io"
	"os"
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

	t.Run("ignores nodes that can't render", func(t *testing.T) {
		n := Div(JoinAttrs("class", Class("party"), ID("hey"), &brokenNode{first: true}, Class("hat")))
		assert.Equal(t, `<div class="party hat" id="hey"></div>`, n)
	})
}

func myButton(children ...g.Node) g.Node {
	return Div(JoinAttrs("class", g.Group(children), Class("button")))
}

func myPrimaryButton(text string) g.Node {
	return myButton(Class("primary"), g.Text(text))
}

func ExampleJoinAttrs() {
	danceButton := myPrimaryButton("Dance")
	_ = danceButton.Render(os.Stdout)
	// Output: <div class="primary button">Dance</div>
}
