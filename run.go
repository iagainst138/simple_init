package sinit

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Run(configPath string) error {
	config, err := Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer stop()

	for _, task := range config.Pre {
		if err := task.Run(ctx, false); err != nil {
			return fmt.Errorf("pre task failed: %w", err)
		}
	}

	var wg sync.WaitGroup
	for _, task := range config.Services {
		wg.Add(1)
		go func(t *Task) { t.Run(ctx, true); wg.Done() }(task)
	}

	wg.Wait()

	return nil
}
