package middleware

import (
	"net/http"
	"time"

	"github.com/matscus/ammunition/monitoring"
	log "github.com/sirupsen/logrus"
)

func Middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		monitoring.RequestCount.Inc()
		start := time.Now()
		defer func() {
			if r := recover(); r != nil {
				log.Debug("Recovered")
			}
			monitoring.RequestsDuration.WithLabelValues(r.Method, r.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
		}()
	}
}
