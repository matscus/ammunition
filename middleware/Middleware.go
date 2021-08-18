package middleware

import (
	"net/http"
	"time"

	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

func Middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics.RequestCount.Inc()
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				log.Error("Middleware recover panic ", err)
			}
			metrics.RequestsDuration.WithLabelValues(r.Method, r.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
		}()
		f(w, r)
	}
}
