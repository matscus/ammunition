package cache

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/monitoring"
	log "github.com/sirupsen/logrus"
)

var (
	CacheMap sync.Map
	ChanMap  sync.Map
)

type Cache struct {
	Name         string
	BigCache     *bigcache.BigCache
	BufferLen    int
	WorkersCount int
	CH           chan string
}

func init() {
	go getCacheMetrics()
}

func CreateDefaultCache(name string, bufferLen int, workers int) (cache Cache, err error) {
	return createDefaultCache(name, bufferLen, workers)
}

func createDefaultCache(name string, bufferLen int, workers int) (cache Cache, err error) {
	config := bigcache.DefaultConfig(5 * time.Minute)
	config.CleanWindow = 0 * time.Minute
	cache.Name = name
	cache.BigCache, err = bigcache.NewBigCache(config)
	cache.BufferLen = bufferLen
	cache.WorkersCount = workers
	cache.CH = make(chan string, cache.BufferLen)
	defer cache.BigCache.Close()
	return cache, err
}

func (c Cache) Init(data []string) {
	CacheMap.Store(c.Name, c)
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
	ch := make(chan string, c.BufferLen)
	for i := 0; i < c.WorkersCount; i++ {
		go c.RunWorker()
	}
	ChanMap.Store(c.Name, ch)
}

func (c Cache) ReInit(data []string) (err error) {
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

func (c Cache) AddValues(data []string) {
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
}

func (c Cache) Delete() error {
	close(c.CH)
	CacheMap.Delete(c.Name)
	return c.BigCache.Reset()
}

func GetCache(name string) (Cache, error) {
	cache, ok := CacheMap.Load(name)
	if ok {
		return cache.(Cache), nil
	}
	return cache.(Cache), errors.New("Cache not found")
}

func (c Cache) RunWorker() {
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

func getCacheMetrics() {
	defer func() {
		recover()
	}()
	var i float64
	for {
		CacheMap.Range(func(k, v interface{}) bool {
			monitoring.CacheLen.WithLabelValues(k.(string)).Set(float64(v.(Cache).BigCache.Len()))
			i++
			return true
		})
		monitoring.CacheCount.Set(i)
		i = 0
		time.Sleep(60 * time.Second)
	}
}
