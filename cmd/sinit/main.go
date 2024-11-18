package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"sinit"
)

func main() {
	config := ""

	flag.StringVar(&config, "config", config, "config file to use (required)")
	flag.Parse()

	if config == "" {
		fmt.Fprintln(os.Stderr, "ERROR: no config file specified.\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	processes := []*sinit.Process{}
	if _, err := os.Stat(config); !os.IsNotExist(err) {
		if data, err := os.ReadFile(config); err == nil {
			if err = json.Unmarshal(data, &processes); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal(err)
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
