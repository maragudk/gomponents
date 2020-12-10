package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	page := Page("Hi!", r.URL.Path)
	_ = page.Render(w)
}

func Page(title, currentPath string) g.Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			StyleEl(Type("text/css"), g.Raw(".is-active{ font-weight: bold }")),
		},
		Body: []g.Node{
			Navbar(currentPath),
			H1(title),
			P(g.Textf("Welcome to the page at %v.", currentPath)),
		},
	})
}

func Navbar(currentPath string) g.Node {
	return Nav(
		NavbarLink("/", "Home", currentPath),
		NavbarLink("/about", "About", currentPath),
	)
}

func NavbarLink(href, name, currentPath string) g.Node {
	return A(Href(href), Classes{"is-active": currentPath == href}, g.Text(name))
}
