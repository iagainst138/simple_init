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

	for _, process := range config.Pre {
		if err := process.Run(false); err != nil {
			return fmt.Errorf("pre task failed: %w", err)
		}
	}

	var wg sync.WaitGroup
	for _, process := range config.Services {
		wg.Add(1)
		go func(p *Process) { p.Run(true); wg.Done() }(process)
	}

	wg.Wait()

	return nil
}
