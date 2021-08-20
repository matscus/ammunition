package cache

import (
	"errors"
	"strconv"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
)

func (c Cache) PersistedInit(data []string) (err error) {
	config := config.DefaultConfig
	config.LifeWindow = c.Life
	config.CleanWindow = c.Clean
	c.BigCache, err = bigcache.NewBigCache(config)
	if err != nil {
		return err
	}
	PersistedCacheMap.Store(c.Name, c)
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
	c.CH = make(chan string, c.BufferLen)
	ChanMap.Store(c.Name, c.CH)
	for i := 0; i < c.WorkersCount; i++ {
		go c.RunWorker()
	}
	return nil
}

func (c Cache) PersistedDelete() error {
	close(c.CH)
	PersistedCacheMap.Delete(c.Name)
	ChanMap.Delete(c.CH)
	return c.BigCache.Close()
}

func CheckPersistedCache(name string) bool {
	_, ok := PersistedCacheMap.Load(name)
	if ok {
		return true
	}
	return false
}

func GetPersistedCache(name string) (Cache, error) {
	cache, ok := PersistedCacheMap.Load(name)
	if ok {
		return cache.(Cache), nil
	}
	return cache.(Cache), errors.New("Cache not found")
}
