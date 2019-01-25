package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
	"github.com/thoas/stats"

	_ "github.com/mattn/go-sqlite3"
)

// Stats ...
var Stats *stats.Stats = stats.New()
var db *sql.DB

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

	if config.StatsDBFreq > 0 {
		db, err := sql.Open("sqlite3", "./stats.db")
		if err != nil {
			log.Fatal("cannot create database file for storing stats with error", err)
		}

		stmt, err := db.Prepare(`create table if not exists stats (
			ts                        int     default (current_timestamp) primary key,
			uptime_sec                real,
			total_status_code_count   text,
			total_count               integer,
			total_response_time_sec   real,
			average_response_time_sec real
		);`)

		if err != nil {
			log.Fatal("cannot create table with error", err)
		}

		stmt.Exec()

		go startStatsDBSavePolling(config.StatsDBFreq, db)
	}

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

	// custom 404 page, also 500 could be added here
	// user would add error html templates to yaml file
	/*r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404"))
		return
	})*/

	if config.Console {
		r.Get("/_console/stats", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Stats.Data())
			return
		})

		r.Get("/_console/config", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(config)
			return
		})

		r.Get("/_console/log", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("log"))
		})

		box := packr.New("assets", "./assets")
		html, _ := box.FindString("console.html")

		r.Get("/_console/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(html))
			return
		})
	}

	// server handler
	server := &http.Server{
		Handler:      r,
		Addr:         config.Listen,
		WriteTimeout: time.Duration(config.Timeout.Write) * time.Second,
		ReadTimeout:  time.Duration(config.Timeout.Read) * time.Second,
	}

	log.Println("==============================")
	log.Printf("listening on http://0.0.0.0%s\n", config.Listen)
	log.Fatal(server.ListenAndServe())

}
