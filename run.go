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

	wg := sync.WaitGroup{}
	for _, process := range config.Services {
		wg.Add(1)
		go func(p *Process) { p.Run(); wg.Done() }(process)
	}

	wg.Wait()

	return nil
}
