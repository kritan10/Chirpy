package utils

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type MetricsConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *MetricsConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *MetricsConfig) ResetHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Swap(0)
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
	}
}

func (cfg *MetricsConfig) MetricsHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/html")
		res.WriteHeader(200)
		data := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
		res.Write([]byte(data))
	}
}
