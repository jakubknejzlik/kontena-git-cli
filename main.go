package main

import (
	"os"

	"github.com/jakubknejzlik/kontena-git-cli/cmd"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Kontena Git CLI"
	app.Usage = "..."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		cmd.GridCommand(),
		cmd.StackCommand(),
	}

	app.Run(os.Args)
}
