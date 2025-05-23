package components_test

import (
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

func hat(children ...g.Node) g.Node {
	return Div(MergeClasses(g.Group(children), Class("hat")))
}

func partyHat(children ...g.Node) g.Node {
	return hat(ID("party-hat"), Class("party"), g.Group(children))
}

func TestMergeClasses(t *testing.T) {
	t.Run("merges classes", func(t *testing.T) {
		n := partyHat(g.Text("Yo."))
		assert.Equal(t, `<div id="party-hat" class="party hat">Yo.</div>`, n)
	})
}

func ExampleClasses() {
	e := g.El("div", Classes{"party-hat": true, "boring-hat": false})
	_ = e.Render(os.Stdout)
	// Output: <div class="party-hat"></div>
}
