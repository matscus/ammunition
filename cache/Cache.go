package cache

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/metrics"
	log "github.com/sirupsen/logrus"
)

var (
	CacheMap sync.Map
	ChanMap  sync.Map
	KV       *bigcache.BigCache
	//initKV bool
)

type Cache struct {
	Name         string
	BigCache     *bigcache.BigCache
	BufferLen    int
	WorkersCount int
	CH           chan string
}

func init() {
	initKV()
	go getCacheMetrics()

}

func New(name string, bufferLen int, workers int, life time.Duration, clean time.Duration) (cashe Cache, err error) {
	return createCache(name, bufferLen, workers, life, clean)
}

func createCache(name string, bufferLen int, workers int, life time.Duration, clean time.Duration) (cashe Cache, err error) {
	cashe.Name = name
	config := config.DefaultConfig
	config.LifeWindow = life
	config.CleanWindow = clean
	cashe.BigCache, err = bigcache.NewBigCache(config)
	cashe.BufferLen = bufferLen
	cashe.WorkersCount = workers
	cashe.CH = make(chan string, cashe.BufferLen)
	return cashe, err
}

func (c Cache) Init(data []string) {
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

func (c Cache) AddValues(data []string) {
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
}

func (c Cache) Delete() error {
	close(c.CH)
	CacheMap.Delete(c.Name)
	ChanMap.Delete(c.CH)
	return c.BigCache.Close()
}

func GetCache(name string) (Cache, error) {
	cache, ok := CacheMap.Load(name)
	if ok {
		return cache.(Cache), nil
	}
	return cache.(Cache), errors.New("Cache not found")
}

func CheckCache(name string) bool {
	_, ok := CacheMap.Load(name)
	if ok {
		return true
	}
	return false
}

func GetChan(name string) (ch chan string, err error) {
	tempChan, ok := ChanMap.Load(name)
	if ok {
		return tempChan.(chan string), nil
	}
	return tempChan.(chan string), errors.New("Chan not found")
}

func (c Cache) RunWorker() {
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

func getCacheMetrics() {
	defer func() {
		recover()
	}()
	//time.Sleep(1 * time.Second)
	var i float64
	metrics.CacheCount.WithLabelValues("persist").Set(i)
	for {
		CacheMap.Range(func(k, v interface{}) bool {
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

func initKV() {
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
