package sinit

import (
	"fmt"
	"sync"
)

func Run(configPath string) error {
	processes, err := Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	wg := sync.WaitGroup{}
	for _, pprocess := range processes {
		wg.Add(1)
		go func(p *Process) {
			p.Run()
			wg.Done()
		}(pprocess)
	}

	wg.Wait()

	return nil
}