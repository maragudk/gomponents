package html_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	. "github.com/maragudk/gomponents/html"
)

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

func TestBooleanAttributes(t *testing.T) {
	cases := map[string]func() g.Node{
		"async":     Async,
		"autofocus": AutoFocus,
		"autoplay":  AutoPlay,
		"controls":  Controls,
		"defer":     Defer,
		"disabled":  Disabled,
		"multiple":  Multiple,
		"readonly":  ReadOnly,
		"required":  Required,
		"selected":  Selected,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := g.El("div", fn())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, name), n)
		})
	}
}

func TestSimpleAttributes(t *testing.T) {
	cases := map[string]func(string) g.Node{
		"accept":       Accept,
		"autocomplete": AutoComplete,
		"charset":      Charset,
		"class":        Class,
		"cols":         Cols,
		"content":      Content,
		"form":         FormAttr,
		"height":       Height,
		"href":         Href,
		"id":           ID,
		"lang":         Lang,
		"max":          Max,
		"maxlength":    MaxLength,
		"min":          Min,
		"minlength":    MinLength,
		"name":         Name,
		"pattern":      Pattern,
		"preload":      Preload,
		"placeholder":  Placeholder,
		"rel":          Rel,
		"rows":         Rows,
		"src":          Src,
		"style":        StyleAttr,
		"tabindex":     TabIndex,
		"target":       Target,
		"title":        TitleAttr,
		"type":         Type,
		"value":        Value,
		"width":        Width,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf(`should output %v="hat"`, name), func(t *testing.T) {
			n := g.El("div", fn("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v="hat"></div>`, name), n)
		})
	}
}
