package sinit

import (
	"fmt"
	"sync"
)

func Run(configPath string) error {
	config, err := Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	for _, task := range config.Pre {
		if err := task.Run(false); err != nil {
			return fmt.Errorf("pre task failed: %w", err)
		}
	}

	var wg sync.WaitGroup
	for _, task := range config.Services {
		wg.Add(1)
		go func(t *Task) { t.Run(true); wg.Done() }(task)
	}

	wg.Wait()

	return nil
}
