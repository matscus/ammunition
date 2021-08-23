package cache

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/metrics"
)

var (
	PersistedCacheMap sync.Map
	ChanMap           sync.Map
	KV                *bigcache.BigCache
	CookiesCache      *bigcache.BigCache
)

type Cache struct {
	Name         string
	BigCache     *bigcache.BigCache
	BufferLen    int
	WorkersCount int
	CH           chan string
	Life         time.Duration
	Clean        time.Duration
}

func init() {
	//initKV()
	go getCacheMetrics()

}

func (c Cache) SetValues(data []string) {
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
}

func GetChan(name string) (ch chan string, err error) {
	tempChan, ok := ChanMap.Load(name)
	if ok {
		return tempChan.(chan string), nil
	}
	return tempChan.(chan string), errors.New("Chan not found")
}

func getCacheMetrics() {
	defer func() {
		recover()
	}()
	var i float64
	metrics.CacheCount.WithLabelValues("kv").Set(1)
	for {
		PersistedCacheMap.Range(func(k, v interface{}) bool {
			metrics.CacheLen.WithLabelValues(k.(string)).Set(float64(v.(Cache).BigCache.Len()))
			metrics.CacheCap.WithLabelValues(k.(string)).Set(float64(v.(Cache).BigCache.Len()))
			i++
			return true
		})
		metrics.CacheCount.WithLabelValues("persist").Set(i)
		i = 0
		metrics.CacheLen.WithLabelValues("kv").Set(float64(KV.Len()))
		metrics.CacheCap.WithLabelValues("kv").Set(float64(KV.Len()))
		time.Sleep(10 * time.Second)
	}
}

// func initKV() {
// 	config := config.DefaultConfig
// 	config.LifeWindow = 1 * time.Hour
// 	config.CleanWindow = 1 * time.Second
// 	var err error
// 	KV, err = bigcache.NewBigCache(config)
// 	if err != nil {
// 		log.Panic("Init KV panic ", err)
// 	}
// 	log.Info("KV init completed")
// }
