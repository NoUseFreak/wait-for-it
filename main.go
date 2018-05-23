package main

import (
	"os"

	"github.com/urfave/cli"
)

var cliUi = new(CliUi)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config-file, f",
			Value: "./wait-for-it.yml",
			Usage: "location of the config file",
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

	pluginLoader, _ := NewPluginLoader(wfiDir+"/plugins")
	pluginLoader.LoadAll(config.Services)

	pluginRunner, _ := NewPluginRunner(wfiDir+"/plugins")
	pluginRunner.RunAll(config.Services)



	//pluginLoader.CleanUp()

	return nil
}
