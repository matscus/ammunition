package cache

import (
	"strconv"
	"time"

	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

func (c Cache) RunWorker() {
	defer func() {
		recover()
	}()
	metrics.WorkerCount.WithLabelValues(c.Name).Inc()
	for {
		for i := 0; i < c.BigCache.Len(); i++ {
			start := time.Now()
			d, err := c.BigCache.Get(strconv.Itoa(i))
			if err != nil {
				log.Println("Worker get values error: ", err)
			}
			c.CH <- string(d)
			metrics.WorkerDuration.WithLabelValues(c.Name).Observe(float64(time.Since(start).Milliseconds()))
		}
	}
}
