package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Process struct {
	CmdPath          string
	Args             []string
	RestartOnFailure bool
	KeepAlive        bool
	RestartWait      int
	WorkingDir       string
	DelayStart       int
}

func (p *Process) Run() {
	log.Printf("starting monitoring %v with args %v", p.CmdPath, p.Args)
	run := true
	for run {
		time.Sleep(time.Duration(p.DelayStart) * time.Second)
		cmd := exec.Command(p.CmdPath, p.Args...)
		if p.WorkingDir != "" {
			cmd.Dir = p.WorkingDir
		}

		r, w, e := os.Pipe()
		if e != nil {
			panic(e)
		}

		cmd.Stdout, cmd.Stderr = w, w
		defer r.Close()
		defer w.Close()

		go func() {
			s := bufio.NewScanner(r)
			for s.Scan() {
				log.Printf("[%v] %v", p.CmdPath, s.Text())
			}
		}()

		err := cmd.Start()
		log.Printf("started %v with args %v", p.CmdPath, p.Args)
		if err != nil {
			log.Printf("error '%v': %v", p.CmdPath, err)
			break // TODO improve on error handling
		}

		err = cmd.Wait()
		r.Close()
		w.Close()
		if err != nil {
			log.Printf("process '%v' finished with error: %v", p.CmdPath, err)
			if !p.RestartOnFailure {
				break // TODO improve
			}
		} else {
			log.Printf("process '%v' finished cleanly", p.CmdPath)
		}

		run = p.KeepAlive
		if run {
			time.Sleep(time.Duration(p.RestartWait) * time.Second)
		}
	}
	log.Printf("finished monitoring '%v'", p.CmdPath)
}

func main() {
	config := ""

	flag.StringVar(&config, "config", config, "config file to use")
	flag.Parse()

	if config == "" {
		log.Fatal("ERROR: no config file specified")
	}

	processes := []*Process{}
	if _, err := os.Stat(config); !os.IsNotExist(err) {
		if data, err := ioutil.ReadFile(config); err == nil {
			if err = json.Unmarshal(data, &processes); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for _, p := range processes {
		wg.Add(1)
		go func(x *Process) {
			x.Run()
			wg.Done()
		}(p)
	}

	wg.Wait()
}
