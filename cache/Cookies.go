package cache

import (
	"encoding/json"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

var (
	cookiesCashe *bigcache.BigCache
	cookiesChan  chan []byte
)

func init() {
	initCookiesCashe()
}

func initCookiesCashe() {
	config := config.DefaultConfig
	config.LifeWindow = 5 * time.Hour
	config.CleanWindow = 1 * time.Second
	var err error
	cookiesCashe, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Panic("Init Cookies panic ", err)
	}
	cookiesChan = make(chan []byte, 1000)
	go cookiesWorker()
	log.Info("Cookies init completed")
}

func SetCookies(key string, values string) error {
	return cookiesCashe.Set(key, []byte(values))
}

func GetCookies() []byte {
	return <-cookiesChan
}

func cookiesWorker() {
	for {
		iterator := cookiesCashe.Iterator()
		start := time.Now()
		for iterator.SetNext() {
			entry, err := iterator.Value()
			if err != nil {
				log.Println(err)
			}
			data := Data{Key: entry.Key(), Value: string(entry.Value())}
			bytes, err := json.Marshal(data)
			if err != nil {
				log.Error("Cookies worker error ", err)
			}
			if len(cookiesChan) < cookiesCashe.Len() {
				cookiesChan <- bytes
				metrics.WorkerDuration.WithLabelValues("cookies").Observe(float64(time.Since(start).Milliseconds()))
			}
		}
	}
}
func getCookiesCacheMetrics() {
	defer func() {
		recover()
	}()
	var i float64
	for {
		metrics.CacheLen.WithLabelValues("cookies").Set(float64(cookiesCashe.Len()))
		metrics.CacheCap.WithLabelValues("cookies").Set(float64(cookiesCashe.Capacity()))
		metrics.CacheCount.WithLabelValues("persist").Set(i)
		i = 0
		time.Sleep(10 * time.Second)
	}
}
