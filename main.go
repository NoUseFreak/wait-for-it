package main

import (
	"os"

	"fmt"
	"github.com/urfave/cli"
)

var cliUi = new(CliUi)

func main() {
	app := cli.NewApp()
	app.Name = "wait-for-it"
	app.Usage = "Wait for services to allow connections"
	app.HideVersion = true
	app.Copyright = "(c) Dries De Peuter <dries@depeuter.io>"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-file, f",
			Value: "./wait-for-it.yml",
			Usage: "location of the config file",
		},
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "toggle to disable cache",
		},
	}

	app.Action = RunAction

	err := app.Run(os.Args)
	if err != nil {
		cliUi.Error(err.Error())
	}
}

func RunAction(c *cli.Context) error {
	wfiDir := "./.wait-for-it"
	config, _ := NewConfig(c.String("config-file"))

	pluginLoader, _ := NewPluginLoader(wfiDir + "/plugins")
	pluginLoader.LoadAll(config.Services, c.Bool("no-cache"))

	pluginRunner, _ := NewPluginRunner(wfiDir + "/plugins")
	completed := pluginRunner.RunAll(config.Services)
	completedLen := len(config.Services)

	cliUi.Title("Report")
	cliUi.Output(fmt.Sprintf("Completed %d/%d", completed, completedLen))
	cliUi.Output(fmt.Sprintf("Failed    %d/%d\n", completedLen-completed, completedLen))

	if completed == completedLen {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

	return nil
}
