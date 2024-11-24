package sinit

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Tasks []*Task

type Config struct {
	Pre      Tasks `yaml:"pre"`
	Services Tasks `yaml:"services"`
}

func (c Config) allTasks() Tasks {
	tasks := make(Tasks, 0, 20)

	tasks = append(tasks, c.Pre...)
	tasks = append(tasks, c.Services...)

	return tasks
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

	for _, t := range config.allTasks() {
		if t.Signal == "" {
			t.Signal = defaultSignal
		} else if !isValidSignal(t.Signal) {
			return config, fmt.Errorf("%s unsupported signal %s - please choose from one of %v", t.Name, t.Signal, validSignals)
		}
	}

	return config, nil
}
