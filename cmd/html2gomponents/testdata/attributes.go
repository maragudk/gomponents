package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Component() Node {
	return A(Href("#"), Title("halløj"), Text("Halløj!"))
}
