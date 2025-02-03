package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Swap(0)
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
	}
}

func (cfg *apiConfig) metricsHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		data := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
		res.Write([]byte(data))
	}
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.ServeMux{}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /healthz", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		body := "OK"
		res.Write([]byte(body))
	})

	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler())

	mux.HandleFunc("POST /reset", apiCfg.resetHandler())

	server := http.Server{Addr: ":8080", Handler: &mux}
	server.ListenAndServe()
}
