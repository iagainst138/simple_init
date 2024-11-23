package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"sinit"
)

func main() {
	configPath := ""

	flag.StringVar(&configPath, "config", configPath, "path to config file to use (required)")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintln(os.Stderr, "ERROR: no config file specified.\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	processes, err := sinit.Load(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	for _, pprocess := range processes {
		wg.Add(1)
		go func(p *sinit.Process) {
			p.Run()
			wg.Done()
		}(pprocess)
	}

	wg.Wait()
}
