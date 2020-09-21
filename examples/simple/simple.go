package main

import (
	"fmt"
	"net/http"
	"time"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
)

func main() {
	_ = http.ListenAndServe("localhost:8080", handler())
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		props := pageProps{
			title: r.URL.Path,
			path:  r.URL.Path,
		}

		_ = g.Write(w, page(props))
	}
}

type pageProps struct {
	title string
	path  string
}

func page(props pageProps) g.Node {
	return el.Document(
		el.HTML(
			g.Attr("lang", "en"),
			el.Head(
				el.Title(props.title),
				el.Style(g.Attr("type", "text/css"), g.Raw(".is-active{font-weight: bold}")),
			),
			el.Body(
				navbar(navbarProps{path: props.path}),
				el.H1(props.title),
				el.P(g.Text(fmt.Sprintf("Welcome to the page at %v.", props.path))),
				el.P(g.Text(fmt.Sprintf("Rendered at %v", time.Now()))),
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
		lis = append(lis, el.Li(
			el.A(item.path, attr.Classes(map[string]bool{"is-active": props.path == item.path}), g.Text(item.text))))
	}
	return el.Ul(lis...)
}
