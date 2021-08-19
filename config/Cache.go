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

type Config struct {
	Cache struct {
		Shards             int  `yaml:"Shards"`
		MaxEntriesInWindow int  `yaml:"MaxEntriesInWindow"`
		MaxEntrySize       int  `yaml:"MaxEntrySize"`
		Verbose            bool `yaml:"Verbose"`
		HardMaxCacheSize   int  `yaml:"HardMaxCacheSize"`
	} `yaml:"Cache"`
}

func init() {
	initConfig()
}

func initConfig() {
	yml, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Panic("Read config file error: ", err)
	}
	config := Config{}
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		log.Panic("Unmarshal config file error: ", err)
	}
	DefaultConfig = bigcache.DefaultConfig(24 * time.Hour)
	DefaultConfig.Shards = config.Cache.Shards
	DefaultConfig.MaxEntriesInWindow = config.Cache.MaxEntriesInWindow
	DefaultConfig.MaxEntrySize = config.Cache.MaxEntrySize
	DefaultConfig.Verbose = config.Cache.Verbose
	DefaultConfig.HardMaxCacheSize = config.Cache.HardMaxCacheSize
}
