package main

import (
	"github.com/mitchellh/colorstring"
	"fmt"
)

type CliUi struct {
}


func (u CliUi) Output(msg string) {
	colorstring.Println(msg)
}

func (u CliUi) Title(msg string) {
	colorstring.Println(fmt.Sprintf("\n[bold][blue]%s", msg))
}

func (u CliUi) Info(msg string) {
	colorstring.Println(fmt.Sprintf("[green]%s", msg))
}

func (u CliUi) Warn(msg string) {
	colorstring.Println(fmt.Sprintf("[blue]%s", msg))
}

func (u CliUi) Error(msg string) {
	colorstring.Println(fmt.Sprintf("[red]%s", msg))
}