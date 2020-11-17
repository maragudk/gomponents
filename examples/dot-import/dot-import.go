package main

import (
	"net/http"

	. "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/el"
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

func page(p props) Node {
	return HTML5(DocumentProps{
		Title:    p.title,
		Language: "en",
		Body: []Node{
			H1(p.title),
			P(Textf("Welcome to the page at %v.", p.path)),
		},
	})
}
