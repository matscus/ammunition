package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ActuatorInfo struct {
	Build struct {
		Application struct {
			Name    string `yaml:"name"`
			Type    string `yaml:"type"`
			Version string `yaml:"version"`
		} `yaml:"application"`
		Author string `yaml:"author"`
	} `yaml:"build"`
}

var Info ActuatorInfo

func ReadActuatorConfig(path string) error {
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yml, &Info)
}
