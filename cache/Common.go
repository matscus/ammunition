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

type Cookies struct {
	Key      string     `json:"key,omitempty"`
	Platform IAMCookies `json:"platform"`
	UFS      string     `json:"ufs"`
}

type IAMCookies struct {
	PlatformSession  string `json:"platformSession,omitempty"`
	PlatformSession2 string `json:"platformSession2,omitempty"`
}
