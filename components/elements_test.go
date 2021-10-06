package components_test

import (
	"testing"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/internal/assert"
)

func TestInputHidden(t *testing.T) {
	t.Run("returns an input element with type hidden, and the given name and value", func(t *testing.T) {
		n := c.InputHidden("id", "partyhat", g.Attr("class", "hat"))
		assert.Equal(t, `<input type="hidden" name="id" value="partyhat" class="hat">`, n)
	})
}

func TestLinkStylesheet(t *testing.T) {
	t.Run("returns a link element with rel stylesheet and the given href", func(t *testing.T) {
		n := c.LinkStylesheet("style.css", g.Attr("media", "print"))
		assert.Equal(t, `<link rel="stylesheet" href="style.css" media="print">`, n)
	})
}

func TestLinkPreload(t *testing.T) {
	t.Run("returns a link element with rel preload and the given href and as", func(t *testing.T) {
		n := c.LinkPreload("party.woff2", "font", g.Attr("type", "font/woff2"))
		assert.Equal(t, `<link rel="preload" href="party.woff2" as="font" type="font/woff2">`, n)
	})
}
