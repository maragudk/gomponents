package el_test

import (
	"errors"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

type erroringWriter struct{}

func (w *erroringWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("don't want to write")
}

func TestDocument(t *testing.T) {
	t.Run("returns doctype and children", func(t *testing.T) {
		assert.Equal(t, `<!doctype html><html></html>`, el.Document(g.El("html")))
	})

	t.Run("errors on write error in Render", func(t *testing.T) {
		err := el.Document(g.El("html")).Render(&erroringWriter{})
		assert.Error(t, err)
	})
}

func TestForm(t *testing.T) {
	t.Run("returns a form element with action and method attributes", func(t *testing.T) {
		assert.Equal(t, `<form action="/" method="post"></form>`, el.Form("/", "post"))
	})
}

func TestInput(t *testing.T) {
	t.Run("returns an input element with attributes type and name", func(t *testing.T) {
		assert.Equal(t, `<input type="text" name="hat">`, el.Input("text", "hat"))
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
		assert.Equal(t, `<progress value="5.5" max="10"></progress>`, el.Progress(5.5, 10))
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
		assert.Equal(t, `<textarea name="hat"></textarea>`, el.Textarea("hat"))
	})
}

func TestA(t *testing.T) {
	t.Run("returns an a element with a href attribute", func(t *testing.T) {
		assert.Equal(t, `<a href="#">hat</a>`, el.A("#", g.Text("hat")))
	})
}

func TestImg(t *testing.T) {
	t.Run("returns an img element with href and alt attributes", func(t *testing.T) {
		assert.Equal(t, `<img src="hat.png" alt="hat" id="image">`, el.Img("hat.png", "hat", g.Attr("id", "image")))
	})
}
