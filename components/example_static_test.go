package components_test

import (
	"os"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// head is the cache slot, a package-level variable next to the component that uses it.
var head string

func ExampleStatic() {
	e := HTML(
		Static(&head, func() g.Node {
			return Head(TitleEl(g.Text("My site")), Link(Rel("stylesheet"), Href("/app.css")))
		}),
		Body(g.Text("Hello")),
	)
	_ = e.Render(os.Stdout)
	// Output: <html><head><title>My site</title><link rel="stylesheet" href="/app.css"></head><body>Hello</body></html>
}
