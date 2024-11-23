package sinit

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	Name             string   `yaml:"name"`
	CmdPath          string   `yaml:"cmd_path"`
	Args             []string `yaml:"args"`
	RestartOnFailure bool     `yaml:"restart_on_failure"`
	Keepalive        bool     `yaml:"keepalive"`
	RestartWait      int      `yaml:"restart_wait"`
	WorkingDir       string   `yaml:"working_dir"`
	DelayStart       int      `yaml:"delay_start"`
}

func (p Process) command() string {
	return fmt.Sprintf("%s %v", p.CmdPath, p.Args)
}

func (p *Process) Run(logInRealtime bool) error {
	var out bytes.Buffer
	var ok bool

	if !logInRealtime {
		defer func() {
			if !ok {
				log.Printf("pre task %q failed:", p.Name)
				fmt.Println(out.String())
			}
		}()
	}

	log.Printf("starting monitoring %q - %s", p.Name, p.command())

	run := true
	for run {
		time.Sleep(time.Duration(p.DelayStart) * time.Second)
		cmd := exec.Command(p.CmdPath, p.Args...)
		if p.WorkingDir != "" {
			cmd.Dir = p.WorkingDir
		}

		r, w, err := os.Pipe()
		if err != nil {
			return err
		}

		cmd.Stdout, cmd.Stderr = w, w
		defer r.Close()
		defer w.Close()

		go func() {
			s := bufio.NewScanner(r)
			for s.Scan() {
				line := fmt.Sprintf("[%v] %v", p.CmdPath, s.Text())
				if logInRealtime {
					log.Println(line)
				} else {
					out.WriteString(line + "\n")
				}
			}
		}()

		log.Printf("starting %q - %s", p.Name, p.command())
		if err = cmd.Start(); err != nil {
			return fmt.Errorf("[%v]: %w", p.CmdPath, err)
		}

		err = cmd.Wait()
		r.Close()
		w.Close()
		if err != nil {
			log.Printf("process %q finished with error: %v", p.Name, err)
			if !p.RestartOnFailure {
				return err
			}
		} else {
			log.Printf("process %q finished cleanly", p.Name)
		}

		run = p.Keepalive
		if run {
			time.Sleep(time.Duration(p.RestartWait) * time.Second)
		}
	}
	log.Printf("finished monitoring %q", p.Name)

	ok = true

	return nil
}
