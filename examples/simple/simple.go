//go:build go1.18
// +build go1.18

package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_ = Page(props{
		title: r.URL.Path,
		path:  r.URL.Path,
	}).Render(w)
}

type props struct {
	title string
	path  string
}

// Page is a whole document to output.
func Page(p props) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    p.title,
		Language: "en",
		Head: []g.Node{
			StyleEl(Type("text/css"),
				g.Raw("html { font-family: sans-serif; }"),
				g.Raw("ul { list-style-type: none; margin: 0; padding: 0; overflow: hidden; }"),
				g.Raw("ul li { display: block; padding: 8px; float: left; }"),
				g.Raw(".is-active { font-weight: bold; }"),
			),
		},
		Body: []g.Node{
			Navbar(p.path, []PageLink{
				{Path: "/foo", Name: "Foo"},
				{Path: "/bar", Name: "Bar"},
			}),
			H1(g.Text(p.title)),
			P(g.Textf("Welcome to the page at %v.", p.path)),
		},
	})
}

type PageLink struct {
	Path string
	Name string
}

func Navbar(currentPath string, links []PageLink) g.Node {
	return Div(
		Ul(
			NavbarLink("/", "Home", currentPath),

			g.Group(g.Map(links, func(pl PageLink) g.Node {
				return NavbarLink(pl.Path, pl.Name, currentPath)
			})),
		),

		Hr(),
	)
}

func NavbarLink(href, name, currentPath string) g.Node {
	return Li(A(Href(href), c.Classes{"is-active": currentPath == href}, g.Text(name)))
}
