package htmx

import g "github.com/maragudk/gomponents"

func Hyper(s string) g.Node {
	return HxAttr(`_`, s)
}
