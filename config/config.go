package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port            int `yaml:"port"`
	ShutdownTimeout int `yaml:"shutdown_timeout"`
}

// MustInit read config file and parse it into struct.
// Panics if any operations fail.
func MustInit(filePath string) *Config {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error while reading config file %q: %v", filePath, err))
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		panic(fmt.Sprintf("Error while unmarshalling configuration: %v", err))
	}
	return cfg
}
