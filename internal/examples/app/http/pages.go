package http

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	ghttp "github.com/maragudk/gomponents/http"

	"app/html"
)

func Home(mux *http.ServeMux) {
	mux.Handle("GET /", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		// Let's pretend this comes from a db or something.
		items := []string{"Super", "Duper", "Nice"}
		return html.HomePage(items), nil
	}))
}

func About(mux *http.ServeMux) {
	mux.Handle("GET /about", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		return html.AboutPage(), nil
	}))
}
