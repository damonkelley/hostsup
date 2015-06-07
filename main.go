package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/damonkelley/hostsup/commands"
)

func main() {
	app := cli.NewApp()
	app.Commands = commands.Commands

	app.Run(os.Args)
}
