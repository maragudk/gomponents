package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePage(items []string) Node {
	return page("Home",
		H1(Text("Home")),

		P(Text("This is a gomponents example app!")),

		P(Raw(`Have a look at the <a href="https:/github.com/maragudk/gomponents/tree/main/internal/examples/app">source code</a> to see how it’s structured.`)),

		Ul(Map(items, func(s string) Node {
			return Li(Text(s))
		})),
	)
}
