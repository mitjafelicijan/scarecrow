package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ServiceRegistry map[string][]string

var serviceRegistry ServiceRegistry

func main() {

	fmt.Println("=============================>")

	serviceRegistry = ServiceRegistry{
		"/v1/*": {"localhost:9010"},
		"/v2/*": {"localhost:9020"},
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Compress(6))

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))

	for path, endpoints := range serviceRegistry {
		for _, endpoint := range endpoints {
			fmt.Println("registring", path, endpoint)
			r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				proxy := httputil.NewSingleHostReverseProxy(&url.URL{
					Scheme: "http",
					Host:   endpoint,
				})
				proxy.ServeHTTP(w, r)
			})
		}
	}

	http.ListenAndServe(":9000", r)
}
