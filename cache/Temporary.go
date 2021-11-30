package cache

import (
	"time"

	"ammunition/config"
	"ammunition/metrics"

	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
)

var (
	temporaryCache *bigcache.BigCache
	temporaryChan  chan []byte
)

func InitTemporary() {
	initTemporary()
	go getTemporaryCacheMetrics()
}

func initTemporary() {
	c := bigcache.DefaultConfig(time.Duration(config.Config.Temporary.LifeWindow) * time.Minute)
	c.CleanWindow = time.Duration(config.Config.Temporary.CleanWindow) * time.Minute
	c.HardMaxCacheSize = config.Config.Temporary.HardMaxCacheSize
	c.MaxEntrySize = config.Config.Temporary.MaxEntrySize
	c.Shards = config.Config.Temporary.Shards
	c.Verbose = config.Config.Temporary.Verbose
	var err error
	temporaryCache, err = bigcache.NewBigCache(c)
	if err != nil {
		log.Panic("Init Cookies panic ", err)
	}
	temporaryChan = make(chan []byte, config.Config.Temporary.BufferLen)
	for i := 0; i < config.Config.Temporary.Worker; i++ {
		go temporaryWorker()
	}
}

func SetTemporaryValue(key string, values []byte) error {
	return temporaryCache.Set(key, values)
}

func GetTemporaryIteratorValue() []byte {
	select {
	case res, ok := <-temporaryChan:
		if ok {
			return res
		} else {
			return []byte("{\"Message\":\"Chan is close\"}")
		}
	default:
		return []byte("{\"Message\":\"Chan is empty\"}")
	}
}
func GetTemporaryValue(key string) ([]byte, error) {
	return temporaryCache.Get(key)
}
func DeleteTemporaryValue(key string) error {
	return temporaryCache.Delete(key)
}

func ResetTemporaryCache() error {
	close(temporaryChan)
	err := temporaryCache.Reset()
	if err != nil {
		return err
	}
	temporaryChan = make(chan []byte, config.Config.Temporary.BufferLen)
	return nil
}

func temporaryWorker() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Temporary worker recover panic ", err)
		}
		go temporaryWorker()
	}()
	for {
		iterator := temporaryCache.Iterator()
		start := time.Now()
		for iterator.SetNext() {
			entry, err := iterator.Value()
			if err != nil {
				log.Error("Worker iterarion ", err)
			} else {
				if len(temporaryChan) < temporaryCache.Len() {
					temporaryChan <- entry.Value()
					metrics.WorkerDuration.WithLabelValues("temporary").Observe(float64(time.Since(start).Milliseconds()))
				}
			}
		}
	}
}
func getTemporaryCacheMetrics() {
	log.Info("Temporary metrics init completed")
	metrics.WorkerCount.WithLabelValues("temporary").Inc()
	metrics.CacheCount.WithLabelValues("in-memory", "temporary").Set(1)
	for {
		metrics.CacheLen.WithLabelValues("in-memory", "temporary").Set(float64(temporaryCache.Len()))
		metrics.CacheCap.WithLabelValues("in-memory", "temporary").Set(float64(temporaryCache.Capacity()))
		time.Sleep(10 * time.Second)
	}
}
