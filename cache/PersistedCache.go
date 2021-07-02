package cache

import (
	"errors"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/monitoring"
	log "github.com/sirupsen/logrus"
)

type PersistedCache struct {
	Name         string
	BigCache     *bigcache.BigCache
	BufferLen    int
	WorkersCount int
	CH           chan string
}

func CreatePersistedCache(name string, bufferLen int, workers int) (cache PersistedCache, err error) {
	return createPersistedCache(name, bufferLen, workers)
}

func createPersistedCache(name string, bufferLen int, workers int) (cache PersistedCache, err error) {
	cache.Name = name
	cache.BigCache, err = bigcache.NewBigCache(config.PersistedCacheConfig)
	cache.BufferLen = bufferLen
	cache.WorkersCount = workers
	cache.CH = make(chan string, cache.BufferLen)
	defer cache.BigCache.Close()
	return cache, err
}

func (c PersistedCache) Init(data []string) {
	CacheMap.Store(c.Name, c)
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
	c.CH = make(chan string, c.BufferLen)
	for i := 0; i < c.WorkersCount; i++ {
		go c.RunWorker()
	}
}

func (c PersistedCache) ReInit(data []string) (err error) {
	close(c.CH)
	err = c.BigCache.Reset()
	if err != nil {
		return err
	}
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
	c.CH = make(chan string, c.BufferLen)
	for i := 0; i < c.WorkersCount; i++ {
		go c.RunWorker()
	}
	return nil
}

func (c PersistedCache) AddValues(data []string) {
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
}

func (c PersistedCache) Delete() error {
	close(c.CH)
	CacheMap.Delete(c.Name)
	return c.BigCache.Reset()
}

func GetPersistedCache(name string) (PersistedCache, error) {
	cache, ok := CacheMap.Load(name)
	if ok {
		return cache.(PersistedCache), nil
	}
	return cache.(PersistedCache), errors.New("Cache not found")
}

func (c PersistedCache) RunWorker() {
	defer func() {
		recover()
	}()
	for {
		for i := 0; i < c.BigCache.Len(); i++ {
			start := time.Now()
			d, err := c.BigCache.Get(strconv.Itoa(i))
			if err != nil {
				log.Println(err)
			}
			c.CH <- string(d)
			monitoring.WorkerDuration.WithLabelValues(c.Name).Observe(float64(time.Since(start).Milliseconds()))
		}
	}
}
