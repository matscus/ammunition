package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

var (
	cookiesCache        *bigcache.BigCache
	cookiesChan         chan []byte
	CookiesMaxCacheSize int
)

func init() {
	initCookiesCache()
	go getCookiesCacheMetrics()
}

func initCookiesCache() {
	config := config.DefaultConfig
	config.LifeWindow = 5 * time.Hour
	config.CleanWindow = 1 * time.Minute
	config.HardMaxCacheSize = CookiesMaxCacheSize
	// config.MaxEntrySize = 500
	config.Shards = 1024
	config.Verbose = false
	var err error
	cookiesCache, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Panic("Init Cookies panic ", err)
	}
	cookiesChan = make(chan []byte, 3000)
	go cookiesWorker()
	log.Info("Cookies init completed")
}

func SetCookies(key string, values []byte) error {
	return cookiesCache.Set(key, values)
}

func GetCookies() []byte {
	select {
	case res, ok := <-cookiesChan:
		if ok {
			return res
		} else {
			return []byte("{\"Message\":\"Chan is close\"}")
		}
	default:
		return []byte("{\"Message\":\"Chan is empty\"}")
	}
}

func ResetCookiesCache() error {
	close(cookiesChan)
	err := cookiesCache.Reset()
	if err != nil {
		return err
	}
	cookiesChan = make(chan []byte, 3000)
	return nil
}

func cookiesWorker() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Cookies worker recover panic ", err)
		}
		go cookiesWorker()
	}()
	for {
		iterator := cookiesCache.Iterator()
		start := time.Now()
		for iterator.SetNext() {
			entry, err := iterator.Value()
			if err != nil {
				log.Error("Worker iterarion ", err)
			} else {
				if len(cookiesChan) < cookiesCache.Len() {
					cookiesChan <- entry.Value()
					metrics.WorkerDuration.WithLabelValues("cookies").Observe(float64(time.Since(start).Milliseconds()))
				}
			}
		}
	}
}
func getCookiesCacheMetrics() {
	log.Info("Cookies metrics init completed")
	metrics.WorkerCount.WithLabelValues("cookies").Inc()
	metrics.CacheCount.WithLabelValues("in-memory", "cookies").Set(1)
	for {
		metrics.CacheLen.WithLabelValues("in-memory", "cookies").Set(float64(cookiesCache.Len()))
		metrics.CacheCap.WithLabelValues("in-memory", "cookies").Set(float64(cookiesCache.Capacity()))
		time.Sleep(10 * time.Second)
	}
}
