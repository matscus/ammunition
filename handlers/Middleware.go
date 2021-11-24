package handlers

import (
	"ammunition/metrics"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "application/json")
		metrics.RequestCount.Inc()
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				log.Error("Middleware recover panic ", err)
			}
		}()
		c.Next()
		metrics.RequestsDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
	}
}
