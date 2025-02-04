package main

import (
	"github/kritan10/Chirpy/utils"
	"net/http"
)

func main() {
	metrics := utils.MetricsConfig{}

	mux := http.ServeMux{}

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", metrics.MiddlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/healthz", utils.HealthzHandler())
	mux.HandleFunc("POST /api/validate_chirp", utils.ValidateChirpHandler())

	mux.HandleFunc("GET /admin/metrics", metrics.MetricsHandler())
	mux.HandleFunc("POST /admin/reset", metrics.ResetHandler())

	server := http.Server{Addr: ":8080", Handler: &mux}
	server.ListenAndServe()
}
