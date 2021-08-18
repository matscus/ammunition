package cache

import (
	"errors"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

type PersistedCache struct {
	Name         string
	BigCache     *bigcache.BigCache
	BufferLen    int
	WorkersCount int
	CH           chan string
}

func CreatePersistedCache(name string, bufferLen int, workers int) (persistedCache PersistedCache, err error) {
	return createPersistedCache(name, bufferLen, workers)
}

func createPersistedCache(name string, bufferLen int, workers int) (persistedCache PersistedCache, err error) {
	persistedCache.Name = name
	persistedCache.BigCache, err = bigcache.NewBigCache(config.PersistedCacheConfig)
	persistedCache.BufferLen = bufferLen
	persistedCache.WorkersCount = workers
	persistedCache.CH = make(chan string, persistedCache.BufferLen)
	return persistedCache, err
}

func (c PersistedCache) Init(data []string) {
	CacheMap.Store(c.Name, c)
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
	c.CH = make(chan string, c.BufferLen)
	ChanMap.Store(c.Name, c.CH)
	for i := 0; i < c.WorkersCount; i++ {
		go c.RunWorker()
	}
}

func (c PersistedCache) AddValues(data []string) {
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
}

func (c PersistedCache) Delete() error {
	close(c.CH)
	CacheMap.Delete(c.Name)
	ChanMap.Delete(c.CH)
	return c.BigCache.Close()
}

func GetPersistedCache(name string) (PersistedCache, error) {
	persistedCache, ok := CacheMap.Load(name)
	if ok {
		return persistedCache.(PersistedCache), nil
	}
	log.Println("Persisted cache not found")
	return persistedCache.(PersistedCache), errors.New("Persisted cache not found")
}

func GetPersistedChan(name string) (ch chan string, err error) {
	tempChan, ok := ChanMap.Load(name)
	if ok {
		return tempChan.(chan string), nil
	}
	return tempChan.(chan string), errors.New("Chan not found")
}

func (c PersistedCache) RunWorker() {
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
