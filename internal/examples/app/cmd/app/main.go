package main

import (
	"log/slog"

	"app/http"
)

func main() {
	if err := http.Start(); err != nil {
		slog.Info("Error starting", "error", err)
	}
}
