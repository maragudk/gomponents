package html

import (
	g "maragu.dev/gomponents"
)

// Aria attributes automatically have their name prefixed with "aria-".
func Aria(name, v string) g.Node {
	return g.Attr("aria-"+name, v)
}

// Data attributes automatically have their name prefixed with "data-".
func Data(name, v string) g.Node {
	return g.Attr("data-"+name, v)
}

// Deprecated: Use [Data] instead.
func DataAttr(name, v string) g.Node {
	return Data(name, v)
}

// Deprecated: Use [Style] instead.
func StyleAttr(v string) g.Node {
	return Style(v)
}

// Deprecated: Use [Title] instead.
func TitleAttr(v string) g.Node {
	return Title(v)
}
