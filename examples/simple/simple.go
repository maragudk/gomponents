package main

import (
	"net/http"
	"time"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := page(props{
		title: r.URL.Path,
		path:  r.URL.Path,
	})
	_ = p.Render(w)
}

type props struct {
	title string
	path  string
}

func page(p props) g.Node {
	return h.Doctype(
		h.HTML(h.Lang("en"),
			h.Head(
				h.TitleEl(p.title),
				h.StyleEl(h.Type("text/css"),
					g.Raw(".is-active{font-weight: bold}"),
					g.Raw("ul.nav { list-style-type: none; margin: 0; padding: 0; overflow: hidden; }"),
					g.Raw("ul.nav li { display: block;  padding: 8px; float: left; }"),
				),
			),
			h.Body(
				navbar(navbarProps{path: p.path}),
				h.Hr(),
				h.H1(p.title),
				h.P(g.Textf("Welcome to the page at %v.", p.path)),
				h.P(g.Textf("Rendered at %v", time.Now())),
			),
		),
	)
}

type navbarProps struct {
	path string
}

func navbar(props navbarProps) g.Node {
	items := []struct {
		path string
		text string
	}{
		{"/", "Home"},
		{"/foo", "Foo"},
		{"/bar", "Bar"},
	}
	lis := g.Map(len(items), func(i int) g.Node {
		item := items[i]
		return h.Li(
			h.A(item.path, c.Classes(map[string]bool{"is-active": props.path == item.path}), g.Text(item.text)),
		)
	})
	return h.Ul(h.Class("nav"), g.Group(lis))
}
