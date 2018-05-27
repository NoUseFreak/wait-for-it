package main

import (
	"fmt"
	"github.com/mitchellh/colorstring"
	"strings"
)

type CliUi struct {
	Quite bool
}

func (u CliUi) Output(msg string) {
	if !u.Quite {
		colorstring.Println(msg)
	}
}

func (u CliUi) Title(msg string) {
	if !u.Quite {
		colorstring.Println(fmt.Sprintf("\n[bold][blue]%s\n%s", msg, strings.Repeat("-", len(msg))))
	}
}

func (u CliUi) Info(msg string) {
	if !u.Quite {
		colorstring.Println(fmt.Sprintf("[green]%s", msg))
	}
}

func (u CliUi) Warn(msg string) {
	if !u.Quite {
		colorstring.Println(fmt.Sprintf("[blue]%s", msg))
	}
}

func (u CliUi) Error(msg string) {
	colorstring.Println(fmt.Sprintf("[red]%s", msg))
}
