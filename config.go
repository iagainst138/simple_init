package sinit

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Processes []*Process

type Config struct {
	Pre      Processes `yaml:"pre"`
	Services Processes `yaml:"services"`
}

func Load(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
