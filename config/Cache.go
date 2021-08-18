package config

import (
	"os"
	"time"

	"github.com/allegro/bigcache"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var (
	PersistedCacheConfig bigcache.Config
	TemporaryCacheConfig bigcache.Config
)

// Config
type Config struct {
	PersistedCache PersistedCache `yaml:"PersistedCache"`
	TemporaryCache TemporaryCache `yaml:"TemporaryCache"`
}

// PersistedCache
type PersistedCache struct {
	Shards             int  `yaml:"Shards"`
	MaxEntriesInWindow int  `yaml:"MaxEntriesInWindow"`
	MaxEntrySize       int  `yaml:"MaxEntrySize"`
	Verbose            bool `yaml:"Verbose"`
	HardMaxCacheSize   int  `yaml:"HardMaxCacheSize"`
}

// TemponaryCache
type TemporaryCache struct {
	LifeWindow         int64 `yaml:"LifeWindow"`
	CleanWindow        int64 `yaml:"CleanWindow"`
	MaxEntriesInWindow int   `yaml:"MaxEntriesInWindow"`
	MaxEntrySize       int   `yaml:"MaxEntrySize"`
	Verbose            bool  `yaml:"Verbose"`
	HardMaxCacheSize   int   `yaml:"HardMaxCacheSize"`
	Shards             int   `yaml:"Shards"`
}

func init() {
	initConfig()
}

func initConfig() {
	yml, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Error("Read config file error: ", err)
		log.Info("Init default values from cache, retention 1h")
		PersistedCacheConfig = bigcache.DefaultConfig(1 * time.Hour)
		PersistedCacheConfig.CleanWindow = 0 * time.Minute
		TemporaryCacheConfig = bigcache.DefaultConfig(1 * time.Hour)
		TemporaryCacheConfig.CleanWindow = 93600 * time.Minute
	}
	config := Config{}
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		log.Error("Unmarshal config file error: ", err)
		log.Info("Init default values from cache, retention temporary cache 1h")
		PersistedCacheConfig = bigcache.DefaultConfig(1 * time.Hour)
		PersistedCacheConfig.CleanWindow = 0 * time.Minute
		TemporaryCacheConfig = bigcache.DefaultConfig(1 * time.Hour)
		TemporaryCacheConfig.CleanWindow = 93600 * time.Minute
	}
	PersistedCacheConfig = bigcache.DefaultConfig(1 * time.Hour)
	PersistedCacheConfig.Shards = config.PersistedCache.Shards
	PersistedCacheConfig.CleanWindow = 0
	PersistedCacheConfig.MaxEntriesInWindow = config.PersistedCache.MaxEntriesInWindow
	PersistedCacheConfig.MaxEntrySize = config.PersistedCache.MaxEntrySize
	PersistedCacheConfig.Verbose = config.PersistedCache.Verbose
	PersistedCacheConfig.HardMaxCacheSize = config.PersistedCache.HardMaxCacheSize

	TemporaryCacheConfig = bigcache.DefaultConfig(time.Duration(config.TemporaryCache.LifeWindow) * time.Second)
	TemporaryCacheConfig.Shards = config.TemporaryCache.Shards
	TemporaryCacheConfig.CleanWindow = time.Duration(config.TemporaryCache.CleanWindow) * time.Second
	TemporaryCacheConfig.MaxEntriesInWindow = config.TemporaryCache.MaxEntriesInWindow
	TemporaryCacheConfig.MaxEntrySize = config.TemporaryCache.MaxEntrySize
	TemporaryCacheConfig.Verbose = config.TemporaryCache.Verbose
	TemporaryCacheConfig.HardMaxCacheSize = config.TemporaryCache.HardMaxCacheSize
}
