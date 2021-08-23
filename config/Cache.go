package config

import (
	"io/ioutil"
	"time"

	"github.com/allegro/bigcache"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var (
	DefaultConfig bigcache.Config
)

type Cache struct {
	DefaultCache struct {
		Shards             int  `yaml:"Shards"`
		MaxEntriesInWindow int  `yaml:"MaxEntriesInWindow"`
		MaxEntrySize       int  `yaml:"MaxEntrySize"`
		Verbose            bool `yaml:"Verbose"`
		HardMaxCacheSize   int  `yaml:"HardMaxCacheSize"`
	} `yaml:"DefaultCache"`
}

func init() {
	initConfig()
}

func initConfig() {
	yml, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Panic("Read config file error: ", err)
	}
	config := Cache{}
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		log.Panic("Unmarshal config file error: ", err)
	}
	DefaultConfig = bigcache.DefaultConfig(24 * time.Hour)
	DefaultConfig.Shards = config.DefaultCache.Shards
	DefaultConfig.MaxEntriesInWindow = config.DefaultCache.MaxEntriesInWindow
	DefaultConfig.MaxEntrySize = config.DefaultCache.MaxEntrySize
	DefaultConfig.Verbose = config.DefaultCache.Verbose
	DefaultConfig.HardMaxCacheSize = config.DefaultCache.HardMaxCacheSize
}
