package cache

import (
	"time"

	"github.com/allegro/bigcache"
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
