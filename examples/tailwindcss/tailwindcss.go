package main

import (
	"net/http"
	"time"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func main() {
	http.Handle("/", createHandler(indexPage()))
	http.Handle("/contact", createHandler(contactPage()))
	http.Handle("/about", createHandler(aboutPage()))

	_ = http.ListenAndServe("localhost:8080", nil)
}

func createHandler(title string, body g.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Rendering a Node is as simple as calling Render and passing an io.Writer
		_ = Page(title, r.URL.Path, body).Render(w)
	}
}

func indexPage() (string, g.Node) {
	return "Welcome!", Div(
		H1(g.Text("Welcome to this example page")),
		P(g.Text("I hope it will make you happy. 😄 It's using TailwindCSS for styling.")),
	)
}

func contactPage() (string, g.Node) {
	return "Contact", Div(
		H1(g.Text("Contact us")),
		P(g.Text("Just do it.")),
	)
}

func aboutPage() (string, g.Node) {
	return "About", Div(
		H1(g.Text("About this site")),
		P(g.Text("This is a site showing off gomponents.")),
	)
}

func Page(title, path string, body g.Node) g.Node {
	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@^2.1.x/dist/base.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@^2.1.x/dist/components.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/@tailwindcss/typography@0.4.x/dist/typography.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@^2.1.x/dist/utilities.min.css")),
		},
		Body: []g.Node{
			Navbar(path, []PageLink{
				{Path: "/contact", Name: "Contact"},
				{Path: "/about", Name: "About"},
			}),
			Container(
				Prose(body),
				PageFooter(),
			),
		},
	})
}

type PageLink struct {
	Path string
	Name string
}

func Navbar(currentPath string, links []PageLink) g.Node {
	return Nav(Class("bg-gray-700 mb-4"),
		Container(
			Div(Class("flex items-center space-x-4 h-16"),
				NavbarLink("/", "Home", currentPath == "/"),

				// We can Map custom slices to Nodes
				g.Group(g.Map(len(links), func(i int) g.Node {
					return NavbarLink(links[i].Path, links[i].Name, currentPath == links[i].Path)
				})),
			),
		),
	)
}

// NavbarLink is a link in the Navbar.
func NavbarLink(path, text string, active bool) g.Node {
	return A(Href(path), g.Text(text),
		// Apply CSS classes conditionally
		c.Classes{
			"px-3 py-2 rounded-md text-sm font-medium focus:outline-none focus:text-white focus:bg-gray-700": true,
			"text-white bg-gray-900":                           active,
			"text-gray-300 hover:text-white hover:bg-gray-700": !active,
		},
	)
}

func Container(children ...g.Node) g.Node {
	return Div(Class("max-w-7xl mx-auto px-2 sm:px-6 lg:px-8"), g.Group(children))
}

func Prose(children ...g.Node) g.Node {
	return Div(Class("prose"), g.Group(children))
}

func PageFooter() g.Node {
	return Footer(Class("prose prose-sm prose-indigo"),
		P(
			// We can use string interpolation directly, like fmt.Sprintf.
			g.Textf("Rendered %v. ", time.Now().Format(time.RFC3339)),

			// Conditional inclusion
			g.If(time.Now().Second()%2 == 0, g.Text("It's an even second.")),
			g.If(time.Now().Second()%2 == 1, g.Text("It's an odd second.")),
		),

		P(A(Href("https://www.gomponents.com"), g.Text("gomponents"))),
	)
}
