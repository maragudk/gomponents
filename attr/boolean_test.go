package attr_test

import (
	"fmt"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/assert"
	"github.com/maragudk/gomponents/attr"
)

func TestBooleanAttributes(t *testing.T) {
	cases := map[string]func() g.Node{
		"async":     attr.Async,
		"autofocus": attr.AutoFocus,
		"autoplay":  attr.AutoPlay,
		"controls":  attr.Controls,
		"defer":     attr.Defer,
		"disabled":  attr.Disabled,
		"multiple":  attr.Multiple,
		"readonly":  attr.ReadOnly,
		"required":  attr.Required,
		"selected":  attr.Selected,
	}

	for name, fn := range cases {
		t.Run(fmt.Sprintf("should output %v", name), func(t *testing.T) {
			n := g.El("div", fn())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, name), n)
		})
	}
}
