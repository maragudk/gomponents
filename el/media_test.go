package el_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/el"
)

func TestImg(t *testing.T) {
	t.Run("returns an img element with href and alt attributes", func(t *testing.T) {
		assert.Equal(t, `<img src="hat.png" alt="hat" id="image" />`, el.Img("hat.png", "hat", g.Attr("id", "image")))
	})
}
