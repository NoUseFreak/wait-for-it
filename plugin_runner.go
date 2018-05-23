package main

import (
	"sync"
	"fmt"
	"os/exec"
	"bufio"
	"time"
	"syscall"
)

type PluginRunner struct {
	location string
}

func NewPluginRunner(location string) (PluginRunner, error) {
	pluginRunner := PluginRunner{
		location:location,
	}

	return pluginRunner, nil
}

func (pr *PluginRunner) RunAll(configs []ServiceConfig) {
	cliUi.Title("Testing services...")
	var wg sync.WaitGroup
	responses := make(chan bool)
	stdout := make(chan string)
	stderr := make(chan string)

	for _, service := range configs {
		wg.Add(1)
		go func(service ServiceConfig) {
			defer wg.Done()

			result := pr.Run(service, stdout, stderr)
			responses <- result
		}(service)
	}


	go func() {
		for logLine := range stdout {
			cliUi.Output(logLine)
		}
	}()
	go func() {
		for logLine := range stderr {
			cliUi.Error(logLine)
		}
	}()
	go func() {
		for response := range responses {
			fmt.Println(response)
		}
	}()

	wg.Wait()
}

func (pr *PluginRunner) Run(config ServiceConfig, stdout chan string, stderr chan string) bool {
	cliUi.Output(fmt.Sprintf(" - Running %s (%s) with %v ", config.Name, config.Type, config.Settings))
	cmd := exec.Command(pr.location+"/"+config.Type)

	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	timer := time.AfterFunc(1 * time.Second, func() {
		cmd.Process.Kill()
		stderr <- "Process timeout"
	})

	if err := cmd.Start(); err != nil {
		stderr <- err.Error()
		return false
	}

	stderrScanner := bufio.NewScanner(stderrPipe)
	go func() {
		for stderrScanner.Scan() {
			stderr <- fmt.Sprintf("  %s:%s", config.Name, stderrScanner.Text())
		}
	}()

	stdoutScanner := bufio.NewScanner(stdoutPipe)
	go func() {
		for stdoutScanner.Scan() {
			stdout <- fmt.Sprintf("  %s:%s", config.Name, stdoutScanner.Text())
		}
	}()

	err := cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				stderr <- fmt.Sprintf("Exit Status: %d", status.ExitStatus())
			}
		}
	}
	timer.Stop()

	return err == nil;
}
