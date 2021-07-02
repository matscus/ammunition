package cache

import (
	"errors"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
)

type TemporaryCache struct {
	Name     string
	BigCache *bigcache.BigCache
}

func CreateTemporaryCache(name string) (cache PersistedCache, err error) {
	return createTemporaryCache(name)
}

func createTemporaryCache(name string) (cache PersistedCache, err error) {
	cache.Name = name
	cache.BigCache, err = bigcache.NewBigCache(config.TemporaryCacheConfig)
	defer cache.BigCache.Close()
	return cache, err
}

func GetTemporaryCache(name string) (TemporaryCache, error) {
	cache, ok := CacheMap.Load(name)
	if ok {
		return cache.(TemporaryCache), nil
	}
	return cache.(TemporaryCache), errors.New("Cache not found")
}

func (t TemporaryCache) Put(k string, v []byte) error {
	return t.BigCache.Set(k, v)
}

func (t TemporaryCache) Get(k string) ([]byte, error) {
	return t.BigCache.Get(k)
}
