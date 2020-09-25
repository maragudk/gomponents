package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestButton(t *testing.T) {
	t.Run("returns a button element", func(t *testing.T) {
		assert.Equal(t, `<button />`, el.Button())
	})
}

func TestForm(t *testing.T) {
	t.Run("returns a form element with action and method attributes", func(t *testing.T) {
		assert.Equal(t, `<form action="/" method="post" />`, el.Form("/", "post"))
	})
}

func TestInput(t *testing.T) {
	t.Run("returns an input element with attributes type and name", func(t *testing.T) {
		assert.Equal(t, `<input type="text" name="hat" />`, el.Input("text", "hat"))
	})
}

func TestLabel(t *testing.T) {
	t.Run("returns a label element with attribute for", func(t *testing.T) {
		assert.Equal(t, `<label for="hat">Hat</label>`, el.Label("hat", g.Text("Hat")))
	})
}

func TestOption(t *testing.T) {
	t.Run("returns an option element with attribute label and content", func(t *testing.T) {
		assert.Equal(t, `<option value="hat">Hat</option>`, el.Option("Hat", "hat"))
	})
}

func TestProgress(t *testing.T) {
	t.Run("returns a progress element with attributes value and max", func(t *testing.T) {
		assert.Equal(t, `<progress value="5.5" max="10" />`, el.Progress(5.5, 10))
	})
}

func TestSelect(t *testing.T) {
	t.Run("returns a select element with attribute name", func(t *testing.T) {
		assert.Equal(t, `<select name="hat"><option value="partyhat">Partyhat</option></select>`,
			el.Select("hat", el.Option("Partyhat", "partyhat")))
	})
}

func TestTextarea(t *testing.T) {
	t.Run("returns a textarea element with attribute name", func(t *testing.T) {
		assert.Equal(t, `<textarea name="hat" />`, el.Textarea("hat"))
	})
}
