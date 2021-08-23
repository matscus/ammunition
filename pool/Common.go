package pool

import (
	"errors"
	"sync"
	"time"

	"github.com/allegro/bigcache"
)

var (
	PersistedCacheMap sync.Map
	ChanMap           sync.Map
	KV                *bigcache.BigCache
	CookiesCache      *bigcache.BigCache
)

type Cache struct {
	Name         string
	BigCache     *bigcache.BigCache
	BufferLen    int
	WorkersCount int
	CH           chan string
	Life         time.Duration
	Clean        time.Duration
}

type Data struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func getChan(name string) (ch chan string, err error) {
	tempChan, ok := ChanMap.Load(name)
	if ok {
		return tempChan.(chan string), nil
	}
	return tempChan.(chan string), errors.New("Chan not found")
}
