package htmx

import (
	"fmt"

	g "github.com/maragudk/gomponents"
)

func WsConnect(v string) g.Node {
	return g.Raw(fmt.Sprintf("ws-connect='%s'", v))
}

func WsSend(v string) g.Node {
	return g.Raw(fmt.Sprintf("ws-send='%s'", v))
}
