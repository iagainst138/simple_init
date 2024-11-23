package sinit

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Processes []*Process

func Load(path string) (Processes, error) {
	var processes Processes

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &processes); err != nil {
		return nil, err
	}

	return processes, nil
}
