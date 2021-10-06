// Package http provides adapters to render gomponents in http handlers.
package http

import (
	"net/http"

	g "github.com/maragudk/gomponents"
)

// Handler is like http.Handler but returns a Node and an error.
// See Adapt for how errors are translated to HTTP responses.
type Handler = func(http.ResponseWriter, *http.Request) (g.Node, error)

type errorWithStatusCode interface {
	StatusCode() int
}

// Adapt a Handler to a http.Handlerfunc.
// The returned Node is rendered to the ResponseWriter, in both normal and error cases.
// If the Handler returns an error, and it implements a "StatusCode() int" method, that HTTP status code is sent
// in the response header. Otherwise, the status code http.StatusInternalServerError (500) is used.
func Adapt(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n, err := h(w, r)
		if err != nil {
			switch v := err.(type) {
			case errorWithStatusCode:
				w.WriteHeader(v.StatusCode())
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		if n == nil {
			return
		}

		if err := n.Render(w); err != nil {
			http.Error(w, "error rendering node: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
