package config

import (
	"github.com/carlosrodriguesf/gobff/tool/logger"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type (
	File struct {
		Version int
		Sources map[string]struct {
			Kind   string    `yaml:"kind"`
			Config yaml.Node `yaml:"config"`
		} `yaml:"sources"`
		Resources map[string]struct {
			CacheDuration time.Duration     `yaml:"cache_duration"`
			Source        string            `yaml:"source"`
			Params        map[string]string `yaml:"params"`
		} `yaml:"resources"`
		Routes []struct {
			Path      string `yaml:"path"`
			Method    string `yaml:"method"`
			Resources map[string]struct {
				As        string   `yaml:"as,omitempty"`
				Output    string   `yaml:"output,omitempty"`
				DependsOn []string `yaml:"depends_on"`
			} `yaml:"resources"`
		} `yaml:"routes"`
	}
)

func Load(logger logger.Logger) (*File, error) {
	logger.Info("opening config.yml")
	file, err := os.Open("config.yml")
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer file.Close()

	config := new(File)
	logger.Info("decoding config.yml")
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		log.Fatalln(err)
	}
	return config, err
}
