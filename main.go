package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var log = logrus.New()

func main() {
	app := cli.NewApp()
	app.Commands = Commands
	app.CommandNotFound = cmdNotFound
	app.Version = Version

	app.Run(os.Args)
}

func cmdNotFound(c *cli.Context, command string) {
	log.Fatalf(
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		c.App.Name,
	)
}
