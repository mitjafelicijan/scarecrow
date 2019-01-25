package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/thoas/stats"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// StatsMiddleware ...
func StatsMiddleware(middleware *stats.Stats) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			beginning, _ := middleware.Begin(w)
			rec := statusRecorder{w, 200}
			next.ServeHTTP(&rec, r)
			middleware.End(beginning, stats.WithStatusCode(rec.status))
		}
		return http.HandlerFunc(fn)
	}
}

func startStatsDBSavePolling(freq int, db *sql.DB) {
	for {
		<-time.After(time.Duration(freq) * time.Second)
		go func() {
			stmt, _ := db.Prepare(`
				insert into stats (uptime_sec, total_status_code_count, total_count, total_response_time_sec, average_response_time_sec)
				values(?,?,?,?,?)
			`)

			totalStatusCodeCount, _ := json.Marshal(Stats.Data().TotalStatusCodeCount)
			_, _ = stmt.Exec(Stats.Data().UpTimeSec, string(totalStatusCodeCount), Stats.Data().TotalCount, Stats.Data().TotalResponseTimeSec, Stats.Data().AverageResponseTimeSec)
		}()
	}
}
