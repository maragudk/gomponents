package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Component() Node {
	return Div(AutoComplete("foo"), AutoFocus(), AutoPlay(), CiteAttr("foo"), ColSpan("foo"), CrossOrigin("foo"), DateTime("foo"), EncType("foo"), FormAttr("foo"), ID("foo"), LabelAttr("foo"), MaxLength("foo"), MinLength("foo"), PlaysInline(), ReadOnly(), RowSpan("foo"), SrcSet("foo"), TabIndex("foo"))
}
