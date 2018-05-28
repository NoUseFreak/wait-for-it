package main

import (
	"bufio"
	"fmt"
	"github.com/NoUseFreak/wait-for-it/plugin"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type PluginRunner struct {
	location string
}

func NewPluginRunner(location string) (PluginRunner, error) {
	pluginRunner := PluginRunner{
		location: location,
	}

	return pluginRunner, nil
}

func (pr *PluginRunner) RunAll(configs []ServiceConfig) int {
	cliUi.Title("Testing services...")
	var wg sync.WaitGroup
	responses := make(chan bool)
	stdout := make(chan string)
	stderr := make(chan string)

	for _, service := range configs {
		cliUi.Output(fmt.Sprintf(" - Running %s (%s)", service.Name, service.Type))
	}
	cliUi.Output("")

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
			if logLine != "" {
				cliUi.Output(logLine)
			}
		}
	}()
	go func() {
		for logLine := range stderr {
			if logLine != "" {
				cliUi.Error(logLine)
			}
		}
	}()
	completed := 0
	go func() {
		var b2i = map[bool]int{false: 0, true: 1}
		for response := range responses {
			completed += b2i[response]
		}
	}()

	wg.Wait()

	return completed
}

func (pr *PluginRunner) Run(config ServiceConfig, stdout chan string, stderr chan string) bool {
	argString := pr.createArguments(config)
	cmd := exec.Command(pr.location+"/"+config.Type, argString)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		stderr <- err.Error()
		return false
	}
	pr.forwardOutput(config, stdoutPipe, stdout)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		stderr <- err.Error()
		return false
	}
	pr.forwardOutput(config, stderrPipe, stderr)

	timer := pr.setTimeoutTimer(config, cmd, stderr)

	if err := cmd.Start(); err != nil {
		stderr <- err.Error()
		return false
	}

	err = cmd.Wait()
	pr.handleResult(err, config, stdout, stderr)
	timer.Stop()

	// force all output to render
	stdout <- ""
	stderr <- ""

	return err == nil
}
func (pr *PluginRunner) handleResult(err error, config ServiceConfig, stdout chan string, stderr chan string) {
	if err == nil {
		stdout <- fmt.Sprintf("%s: Connected", config.Name)
	} else {
		stderr <- fmt.Sprintf("%s: Failed to connect", config.Name)
	}
}

func (pr *PluginRunner) forwardOutput(config ServiceConfig, input io.ReadCloser, output chan string) {
	scanner := bufio.NewScanner(input)
	go func() {
		for scanner.Scan() {
			output <- fmt.Sprintf("%s: %s", config.Name, scanner.Text())
		}
	}()
}

func (pr *PluginRunner) setTimeoutTimer(config ServiceConfig, cmd *exec.Cmd, stderr chan string) *time.Timer {
	return time.AfterFunc(config.Timeout, func() {
		cmd.Process.Kill()
	})
}

func (pr *PluginRunner) createArguments(config ServiceConfig) string {
	vsm := []string{}
	for i, v := range config.Settings {
		if i != "" {
			vsm = append(vsm, fmt.Sprintf("%s=%s", i, v))
		}
	}
	return strings.Join(vsm, plugin.ARG_SEPERATOR)
}
