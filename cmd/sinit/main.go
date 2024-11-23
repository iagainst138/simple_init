package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

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

	if err := sinit.Run(configPath); err != nil {
		slog.Error("sinit run failed", "error", err)
		os.Exit(1)
	}
}
