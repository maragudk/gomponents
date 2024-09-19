package http

import (
	"net/http"
)

func Start() error {
	return http.ListenAndServe(":8080", setupRoutes())
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()

	Home(mux)
	About(mux)

	return mux
}
