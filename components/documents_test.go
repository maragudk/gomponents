package components_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/attr"
	c "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/el"
)

func TestHTML5(t *testing.T) {
	t.Run("returns an html5 document template", func(t *testing.T) {
		e := c.HTML5(c.DocumentProps{
			Title:       "Hat",
			Description: "Love hats.",
			Language:    "en",
			Head:        []g.Node{el.Link(attr.Rel("stylesheet"), attr.Href("/hat.css"))},
			Body:        []g.Node{el.Div()},
		})

		assert.Equal(t, `<!doctype html><html lang="en"><head><meta charset="utf-8" /><meta name="viewport" content="width=device-width, initial-scale=1" /><title>Hat</title><meta name="description" content="Love hats." /><link rel="stylesheet" href="/hat.css" /></head><body><div /></body></html>`, e)
	})

	t.Run("returns no language, description, and extra head/body elements if empty", func(t *testing.T) {
		e := c.HTML5(c.DocumentProps{
			Title: "Hat",
		})

		assert.Equal(t, `<!doctype html><html><head><meta charset="utf-8" /><meta name="viewport" content="width=device-width, initial-scale=1" /><title>Hat</title></head><body /></html>`, e)
	})
}
