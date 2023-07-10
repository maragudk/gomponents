package htmx

import (
	g "github.com/maragudk/gomponents"
)

func SSEConnect(v string) g.Node {
	return HxAttr("sse-connect", v)
}

func SSESwap(v string) g.Node {
	return HxAttr("sse-swap", v)
}
