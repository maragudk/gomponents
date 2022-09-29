package html_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/maragudk/gomponents/internal/assert"
)

func TestBooleanAttributes(t *testing.T) {
	cases := map[string]func() g.Node{
		"async":       Async,
		"autofocus":   AutoFocus,
		"autoplay":    AutoPlay,
		"controls":    Controls,
		"defer":       Defer,
		"disabled":    Disabled,
		"loop":        Loop,
		"multiple":    Multiple,
		"muted":       Muted,
		"playsinline": PlaysInline,
		"readonly":    ReadOnly,
		"required":    Required,
		"selected":    Selected,
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
		"action":       Action,
		"alt":          Alt,
		"as":           As,
		"autocomplete": AutoComplete,
		"charset":      Charset,
		"class":        Class,
		"cols":         Cols,
		"content":      Content,
		"enctype":      EncType,
		"for":          For,
		"form":         FormAttr,
		"height":       Height,
		"href":         Href,
		"id":           ID,
		"lang":         Lang,
		"loading":      Loading,
		"max":          Max,
		"maxlength":    MaxLength,
		"method":       Method,
		"min":          Min,
		"minlength":    MinLength,
		"name":         Name,
		"pattern":      Pattern,
		"placeholder":  Placeholder,
		"poster":       Poster,
		"preload":      Preload,
		"rel":          Rel,
		"role":         Role,
		"rows":         Rows,
		"src":          Src,
		"srcset":       SrcSet,
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

func TestAria(t *testing.T) {
	t.Run("returns an attribute which name is prefixed with aria-", func(t *testing.T) {
		n := Aria("selected", "true")
		assert.Equal(t, ` aria-selected="true"`, n)
	})
}

func TestDataAttr(t *testing.T) {
	t.Run("returns an attribute which name is prefixed with data-", func(t *testing.T) {
		n := DataAttr("id", "partyhat")
		assert.Equal(t, ` data-id="partyhat"`, n)
	})
}
