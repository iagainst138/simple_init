package sinit

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"syscall"
)

const defaultSignal = "term"

var (
	signals = map[string]os.Signal{
		"hup":   syscall.SIGHUP,
		"int":   syscall.SIGINT,
		"term":  syscall.SIGTERM,
		"winch": syscall.SIGWINCH,
	}

	validSignals = slices.Collect(maps.Keys(signals))
)

func isValidSignal(s string) bool {
	_, exists := signals[s]
	return exists
}

func getSignal(s string) os.Signal {
	signal, exists := signals[s]
	if !exists {
		panic(fmt.Errorf("%s is not supported", s))
	}
	return signal
}
