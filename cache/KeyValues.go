package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

var kv *bigcache.BigCache

func init() {
	initKeyValuesCache()
	go keyValuesMetrics()
}

func initKeyValuesCache() {
	config := config.DefaultConfig
	config.LifeWindow = 1 * time.Hour
	config.CleanWindow = 1 * time.Second
	var err error
	kv, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Panic("Init KV panic ", err)
	}
	log.Info("KV init completed")
}

func KVSet(key string, values string) error {
	return kv.Set(key, []byte(values))
}
func KVGet(key string) (string, error) {
	res, err := kv.Get(key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
func KVDelete(key string) error {
	return kv.Delete(key)
}

func keyValuesMetrics() {
	defer func() {
		recover()
	}()
	metrics.CacheCount.WithLabelValues("in-memory", "keyValues").Set(1)
	for {
		metrics.CacheLen.WithLabelValues("in-memory", "keyValues").Set(float64(kv.Len()))
		metrics.CacheCap.WithLabelValues("in-memory", "keyValues").Set(float64(kv.Capacity()))
		time.Sleep(10 * time.Second)
	}
}
