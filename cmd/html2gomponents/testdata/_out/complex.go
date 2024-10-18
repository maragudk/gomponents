package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Component() Node {
	return Div(
		H1(ID("title"), Class("pretty"),
			Text("Halløj!"),
		),
		H2(ID("subtitle"), Class("prettier"),
			Text("What is this?"),
		),
		P(Class("prettiest"),
			Text("It's a parser and converter for converting HTML to gomponents Go code."),
		),
	)
}
