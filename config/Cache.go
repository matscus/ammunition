package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	Config config
)

type config struct {
	Default   data `yaml:"Default"`
	Persist   data `yaml:"Persist"`
	Temporary data `yaml:"Temporary"`

}

type data struct {
	Verbose            bool `yaml:"Verbose"`
	HardMaxCacheSize   int  `yaml:"HardMaxCacheSize"`
	Shards             int  `yaml:"Shards"`
	LifeWindow         int  `yaml:"LifeWindow"`
	CleanWindow        int  `yaml:"CleanWindow"`
	MaxEntriesInWindow int  `yaml:"MaxEntriesInWindow"`
	MaxEntrySize       int  `yaml:"MaxEntrySize"`
	Worker             int  `yaml:"Worker"`
	BufferLen          int  `yaml:"BufferLen"`
}

func ReadCacheConfig(path string) error {
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yml, &Config)
}
