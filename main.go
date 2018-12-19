package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thoas/stats"
)

var Stats *stats.Stats = stats.New()

func main() {

	// https://stackoverflow.com/questions/37321760/how-to-set-up-lets-encrypt-for-a-go-server-application
	//certManager := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("example.com"), //Your domain here
	//	Cache:      autocert.DirCache("certs"),            //Folder for storing certificates
	//}

	// parsing config file
	config = parseConfigFile("scarecrow.yml")

	r := chi.NewRouter()

	if config.Logging {
		log.SetOutput(os.Stdout)
		if config.Verbose {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
		}
		r.Use(middleware.Logger)
	}

	if config.Metrics {
		r.Use(StatsMiddleware(Stats))
	}

	if config.GZIP {
		r.Use(middleware.Compress(6))
	}

	if config.Heartbeat {
		r.Use(middleware.Heartbeat("/_heartbeat"))
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	for _, endpoint := range config.ServiceRegistry {
		log.Println(">> registring", endpoint.Path, endpoint.Proxy)
		path := endpoint.Path
		proxy := endpoint.Proxy
		r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   proxy,
			})
			proxy.ServeHTTP(w, r)
		})

	}

	if config.Console {
		r.Get("/_stats", func(w http.ResponseWriter, r *http.Request) {
			var payload = make(map[string]interface{})
			payload["services"] = config.ServiceRegistry
			payload["stats"] = Stats.Data()

			w.Header().Set("Content-Type", "application/javascript")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
			return
		})
	}

	// server handler
	server := &http.Server{
		Handler:      r,
		Addr:         config.Listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("==============================")
	log.Println("listening on", config.Listen)
	log.Fatal(server.ListenAndServe())

}
