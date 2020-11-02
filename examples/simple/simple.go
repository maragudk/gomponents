package main

import (
	"net/http"
	"time"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
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
	return el.Document(
		el.HTML(attr.Lang("en"),
			el.Head(
				el.Title(p.title),
				el.Style(attr.Type("text/css"),
					g.Raw(".is-active{font-weight: bold}"),
					g.Raw("ul.nav { list-style-type: none; margin: 0; padding: 0; overflow: hidden; }"),
					g.Raw("ul.nav li { display: block;  padding: 8px; float: left; }"),
				),
			),
			el.Body(
				navbar(navbarProps{path: p.path}),
				el.Hr(),
				el.H1(p.title),
				el.P(g.Textf("Welcome to the page at %v.", p.path)),
				el.P(g.Textf("Rendered at %v", time.Now())),
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
	var lis []g.Node
	for _, item := range items {
		lis = append(lis, el.Li(el.A(item.path,
			attr.Classes(map[string]bool{"is-active": props.path == item.path}),
			g.Text(item.text))))
	}
	return el.Ul(attr.Class("nav"), g.Group(lis))
}
