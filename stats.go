package main

import (
	"net/http"

	"github.com/thoas/stats"
)

func StatsMiddleware(middleware *stats.Stats) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			beginning, recorder := middleware.Begin(w)
			next.ServeHTTP(w, r)
			middleware.End(beginning, stats.WithRecorder(recorder))
		}
		return http.HandlerFunc(fn)
	}
}
