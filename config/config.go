package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"time"
)

type (
	File struct {
		Version   int
		Resources map[string]struct {
			CacheDuration time.Duration `yaml:"cache_duration"`
			Source        struct {
				Kind   string    `yaml:"kind"`
				Config yaml.Node `yaml:"config"`
			} `yaml:"source"`
		} `yaml:"resources"`
		Routes []struct {
			Path      string `yaml:"path"`
			Method    string `yaml:"method"`
			Resources map[string]struct {
				As string `yaml:"as"`
			} `yaml:"resources"`
		} `yaml:"routes"`
	}
)

func Load(input io.ReadCloser) (*File, error) {
	config := new(File)
	err := yaml.NewDecoder(input).Decode(config)
	if err != nil {
		log.Fatalln(err)
	}
	return config, err
}
