package sinit

import (
	"encoding/json"
	"os"
)

type Processes []*Process

func Load(path string) (Processes, error) {
	var processes Processes

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &processes); err != nil {
		return nil, err
	}

	return processes, nil
}
