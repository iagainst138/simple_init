package sinit

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	CmdPath          string   `yaml:"cmd_path"`
	Args             []string `yaml:"args"`
	RestartOnFailure bool     `yaml:"restart_on_failure"`
	Keepalive        bool     `yaml:"keepalive"`
	RestartWait      int      `yaml:"restart_wait"`
	WorkingDir       string   `yaml:"working_dir"`
	DelayStart       int      `yaml:"delay_start"`
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

		run = p.Keepalive
		if run {
			time.Sleep(time.Duration(p.RestartWait) * time.Second)
		}
	}
	log.Printf("finished monitoring '%v'", p.CmdPath)
}
