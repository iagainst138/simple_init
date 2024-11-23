package sinit

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Tasks []*Task

type Config struct {
	Pre      Tasks `yaml:"pre"`
	Services Tasks `yaml:"services"`
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
