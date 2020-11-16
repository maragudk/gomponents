package attr_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/attr"
)

func TestSimpleAttributes(t *testing.T) {
	cases := map[string]func(string) g.Node{
		"accept":       attr.Accept,
		"autocomplete": attr.AutoComplete,
		"charset":      attr.Charset,
		"class":        attr.Class,
		"content":      attr.Content,
		"form":         attr.Form,
		"height":       attr.Height,
		"href":         attr.Href,
		"id":           attr.ID,
		"lang":         attr.Lang,
		"max":          attr.Max,
		"maxlength":    attr.MaxLength,
		"min":          attr.Min,
		"minlength":    attr.MinLength,
		"name":         attr.Name,
		"pattern":      attr.Pattern,
		"preload":      attr.Preload,
		"placeholder":  attr.Placeholder,
		"rel":          attr.Rel,
		"src":          attr.Src,
		"style":        attr.Style,
		"tabindex":     attr.TabIndex,
		"target":       attr.Target,
		"title":        attr.Title,
		"type":         attr.Type,
		"value":        attr.Value,
		"width":        attr.Width,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf(`should output %v="hat"`, name), func(t *testing.T) {
			n := g.El("div", fn("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v="hat"></div>`, name), n)
		})
	}
}
