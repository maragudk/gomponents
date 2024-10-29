package html_test

import (
	"fmt"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"maragu.dev/gomponents/internal/assert"
)

func TestBooleanAttributes(t *testing.T) {
	tests := []struct {
		Name string
		Func func() g.Node
	}{
		{Name: "async", Func: Async},
		{Name: "autofocus", Func: AutoFocus},
		{Name: "autoplay", Func: AutoPlay},
		{Name: "checked", Func: Checked},
		{Name: "controls", Func: Controls},
		{Name: "defer", Func: Defer},
		{Name: "disabled", Func: Disabled},
		{Name: "loop", Func: Loop},
		{Name: "multiple", Func: Multiple},
		{Name: "muted", Func: Muted},
		{Name: "playsinline", Func: PlaysInline},
		{Name: "readonly", Func: ReadOnly},
		{Name: "required", Func: Required},
		{Name: "selected", Func: Selected},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			n := g.El("div", test.Func())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, test.Name), n)
		})
	}
}

func TestSimpleAttributes(t *testing.T) {
	tests := []struct {
		Name string
		Func func(string) g.Node
	}{
		{Name: "accept", Func: Accept},
		{Name: "action", Func: Action},
		{Name: "alt", Func: Alt},
		{Name: "as", Func: As},
		{Name: "autocomplete", Func: AutoComplete},
		{Name: "charset", Func: Charset},
		{Name: "cite", Func: CiteAttr},
		{Name: "class", Func: Class},
		{Name: "cols", Func: Cols},
		{Name: "colspan", Func: ColSpan},
		{Name: "content", Func: Content},
		{Name: "crossorigin", Func: CrossOrigin},
		{Name: "datetime", Func: DateTime},
		{Name: "draggable", Func: Draggable},
		{Name: "enctype", Func: EncType},
		{Name: "dir", Func: Dir},
		{Name: "for", Func: For},
		{Name: "form", Func: FormAttr},
		{Name: "height", Func: Height},
		{Name: "hidden", Func: Hidden},
		{Name: "href", Func: Href},
		{Name: "id", Func: ID},
		{Name: "integrity", Func: Integrity},
		{Name: "label", Func: LabelAttr},
		{Name: "lang", Func: Lang},
		{Name: "list", Func: List},
		{Name: "loading", Func: Loading},
		{Name: "max", Func: Max},
		{Name: "maxlength", Func: MaxLength},
		{Name: "method", Func: Method},
		{Name: "min", Func: Min},
		{Name: "minlength", Func: MinLength},
		{Name: "name", Func: Name},
		{Name: "pattern", Func: Pattern},
		{Name: "placeholder", Func: Placeholder},
		{Name: "popovertarget", Func: PopoverTarget},
		{Name: "poster", Func: Poster},
		{Name: "preload", Func: Preload},
		{Name: "rel", Func: Rel},
		{Name: "role", Func: Role},
		{Name: "rows", Func: Rows},
		{Name: "rowspan", Func: RowSpan},
		{Name: "src", Func: Src},
		{Name: "srcset", Func: SrcSet},
		{Name: "step", Func: Step},
		{Name: "style", Func: Style},
		{Name: "style", Func: StyleAttr},
		{Name: "tabindex", Func: TabIndex},
		{Name: "target", Func: Target},
		{Name: "title", Func: Title},
		{Name: "title", Func: TitleAttr},
		{Name: "type", Func: Type},
		{Name: "value", Func: Value},
		{Name: "width", Func: Width},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			n := g.El("div", test.Func("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v="hat"></div>`, test.Name), n)
		})
	}
}

func TestVariadicAttributes(t *testing.T) {
	tests := []struct {
		Name string
		Func func(...string) g.Node
	}{
		{Name: "popover", Func: Popover},
	}

	for _, test := range tests {
		t.Run(test.Name + "(no args)", func(t *testing.T) {
			n := g.El("div", test.Func())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, test.Name), n)
		})

		t.Run(test.Name +"(one arg)", func(t *testing.T) {
			n := g.El("div", test.Func("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v="hat"></div>`, test.Name), n)
		})

		t.Run(test.Name + "(two args panics)", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected a panic")
				}
			}()
			n := g.El("div", test.Func("hat", "party"))
			assert.Equal(t, "unreachable", n)
		})
	}
}

func TestAria(t *testing.T) {
	t.Run("returns an attribute which name is prefixed with aria-", func(t *testing.T) {
		n := Aria("selected", "true")
		assert.Equal(t, ` aria-selected="true"`, n)
	})
}

func TestData(t *testing.T) {
	t.Run("returns an attribute which name is prefixed with data-", func(t *testing.T) {
		n := Data("id", "partyhat")
		assert.Equal(t, ` data-id="partyhat"`, n)

		n = DataAttr("id", "partyhat")
		assert.Equal(t, ` data-id="partyhat"`, n)
	})
}
