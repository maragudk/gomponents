package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func InputHidden(name, value string, children ...g.Node) g.Node {
	return Input(Type("hidden"), Name(name), Value(value), g.Group(children))
}
