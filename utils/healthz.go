package utils

import "net/http"

func HealthzHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		body := "OK"
		res.Write([]byte(body))
	}
}
