package cache

import (
	"time"

	"ammunition/config"
	"ammunition/metrics"

	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
)

var (
	defaultCache *bigcache.BigCache
	defaultChan  chan []byte
)

func InitDefault() {
	initDefault()
	go getDefaultCacheMetrics()
}

func initDefault() {
	c := bigcache.DefaultConfig(time.Duration(config.Config.Default.LifeWindow) * time.Minute)
	c.CleanWindow = time.Duration(config.Config.Default.CleanWindow) * time.Minute
	c.HardMaxCacheSize = config.Config.Default.HardMaxCacheSize
	c.MaxEntrySize = config.Config.Default.MaxEntrySize
	c.Shards = config.Config.Default.Shards
	c.Verbose = config.Config.Default.Verbose
	var err error
	defaultCache, err = bigcache.NewBigCache(c)
	if err != nil {
		log.Panic("Init Cookies panic ", err)
	}
	defaultChan = make(chan []byte, config.Config.Default.BufferLen)
	for i := 0; i < config.Config.Default.Worker; i++ {
		go DefaultWorker()
	}
}

func SetDefaultValue(key string, values []byte) error {
	return defaultCache.Set(key, values)
}

func GetDefaultIteratorValue() []byte {
	select {
	case res, ok := <-defaultChan:
		if ok {
			return res
		} else {
			return []byte("{\"Message\":\"Chan is close\"}")
		}
	default:
		return []byte("{\"Message\":\"Chan is empty\"}")
	}
}
func GetDefaultValue(key string) ([]byte, error) {
	return defaultCache.Get(key)
}
func DeleteDefaultValue(key string) error {
	return defaultCache.Delete(key)
}

func ResetDefaultCache() error {
	close(defaultChan)
	err := defaultCache.Reset()
	if err != nil {
		return err
	}
	defaultChan = make(chan []byte, config.Config.Default.BufferLen)
	return nil
}

func DefaultWorker() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Default worker recover panic ", err)
		}
		go DefaultWorker()
	}()
	for {
		iterator := defaultCache.Iterator()
		start := time.Now()
		for iterator.SetNext() {
			entry, err := iterator.Value()
			if err != nil {
				log.Error("Worker iterarion ", err)
			} else {
				if len(defaultChan) < defaultCache.Len() {
					defaultChan <- entry.Value()
					metrics.WorkerDuration.WithLabelValues("Default").Observe(float64(time.Since(start).Milliseconds()))
				}
			}
		}
	}
}
func getDefaultCacheMetrics() {
	log.Info("Default metrics init completed")
	metrics.WorkerCount.WithLabelValues("default").Inc()
	metrics.CacheCount.WithLabelValues("in-memory", "default").Set(1)
	for {
		metrics.CacheLen.WithLabelValues("in-memory", "default").Set(float64(defaultCache.Len()))
		metrics.CacheCap.WithLabelValues("in-memory", "default").Set(float64(defaultCache.Capacity()))
		time.Sleep(10 * time.Second)
	}
}
