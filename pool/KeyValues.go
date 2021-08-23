package pool

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	log "github.com/sirupsen/logrus"
)

var kv *bigcache.BigCache

func init() {
	initKeyValuesCache()
}

func initKeyValuesCache() {
	config := config.DefaultConfig
	config.LifeWindow = 1 * time.Hour
	config.CleanWindow = 1 * time.Second
	var err error
	kv, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Panic("Init KV panic ", err)
	}
	log.Info("KV init completed")
}
