package sinit

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

type Task struct {
	Name             string   `yaml:"name"`
	CmdPath          string   `yaml:"cmd_path"`
	Args             []string `yaml:"args"`
	RestartOnFailure bool     `yaml:"restart_on_failure"`
	Keepalive        bool     `yaml:"keepalive"`
	RestartWait      int      `yaml:"restart_wait"`
	WorkingDir       string   `yaml:"working_dir"`
	DelayStart       int      `yaml:"delay_start"`
	Signal           string   `yaml:"signal"`
}

func (t Task) command() string {
	return fmt.Sprintf("%s %v", t.CmdPath, t.Args)
}

func (t *Task) Run(ctx context.Context, logInRealtime bool) error {
	var out bytes.Buffer
	var ok bool

	if !logInRealtime {
		defer func() {
			if !ok {
				log.Printf("pre task %q failed:", t.Name)
				fmt.Println(out.String())
			}
		}()
	}

	log.Printf("starting monitoring %q - %s", t.Name, t.command())

	run := true
	for run {
		time.Sleep(time.Duration(t.DelayStart) * time.Second)

		cmd := exec.CommandContext(ctx, t.CmdPath, t.Args...)

		if t.WorkingDir != "" {
			cmd.Dir = t.WorkingDir
		}

		if runtime.GOOS != "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true,
			}
		}

		cmd.Cancel = func() error {
			if cmd.Process != nil {
				return cmd.Process.Signal(getSignal(t.Signal))
			}
			return nil
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
				line := fmt.Sprintf("[%v] %v", t.Name, s.Text())
				if logInRealtime {
					log.Println(line)
				} else {
					out.WriteString(line + "\n")
				}
			}
		}()

		log.Printf("starting %q - %s", t.Name, t.command())
		if err = cmd.Start(); err != nil {
			return fmt.Errorf("[%v]: %w", t.CmdPath, err)
		}

		err = cmd.Wait()
		r.Close()
		w.Close()
		if err != nil {
			log.Printf("task %q finished with error: %v", t.Name, err)
			if !t.RestartOnFailure {
				return err
			}
		} else {
			log.Printf("task %q finished cleanly", t.Name)
		}

		run = t.Keepalive
		if run {
			time.Sleep(time.Duration(t.RestartWait) * time.Second)
		}
	}
	log.Printf("finished monitoring %q", t.Name)

	ok = true

	return nil
}
