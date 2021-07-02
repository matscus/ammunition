package cache

import (
	"sync"
	"time"

	"github.com/matscus/ammunition/monitoring"
)

var (
	CacheMap sync.Map
	ChanMap  sync.Map
)

func init() {
	go getCacheMetrics()
}

func getCacheMetrics() {
	defer func() {
		recover()
	}()
	var i float64
	for {
		CacheMap.Range(func(k, v interface{}) bool {
			monitoring.CacheLen.WithLabelValues(k.(string)).Set(float64(v.(PersistedCache).BigCache.Len()))
			i++
			return true
		})
		monitoring.CacheCount.Set(i)
		i = 0
		time.Sleep(60 * time.Second)
	}
}
