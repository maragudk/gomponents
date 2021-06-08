package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func InputHidden(name, value string, children ...g.Node) g.Node {
	return Input(Type("hidden"), Name(name), Value(value), g.Group(children))
}

func LinkStylesheet(href string, children ...g.Node) g.Node {
	return Link(Rel("stylesheet"), Href(href), g.Group(children))
}

func LinkPreload(href, as string, children ...g.Node) g.Node {
	return Link(Rel("preload"), Href(href), As(as), g.Group(children))
}
