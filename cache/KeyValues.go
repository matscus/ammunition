package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

var KV *bigcache.BigCache

func init() {
	initKeyValuesCache()
	go keyValuesMetrics()
}

func initKeyValuesCache() {
	config := config.DefaultConfig
	config.LifeWindow = 1 * time.Hour
	config.CleanWindow = 1 * time.Second
	var err error
	KV, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Panic("Init KV panic ", err)
	}
	log.Info("KV init completed")
}

func keyValuesMetrics() {
	defer func() {
		recover()
	}()
	metrics.CacheCount.WithLabelValues("keyValues").Set(1)
	for {
		metrics.CacheLen.WithLabelValues("keyValues").Set(float64(KV.Len()))
		metrics.CacheCap.WithLabelValues("keyValues").Set(float64(KV.Capacity()))
		time.Sleep(10 * time.Second)
	}
}
