package http

import (
	"net/http"

	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"

	"app/html"
)

func Home(mux *http.ServeMux) {
	mux.Handle("GET /", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Let's pretend this comes from a db or something.
		items := []string{"Super", "Duper", "Nice"}
		return html.HomePage(items), nil
	}))
}

func About(mux *http.ServeMux) {
	mux.Handle("GET /about", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return html.AboutPage(), nil
	}))
}
