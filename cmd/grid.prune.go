package cmd

import (
	"fmt"
	"os"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/kontena"

	"github.com/urfave/cli"
)

func pruneStacksCommand() cli.Command {
	return cli.Command{
		Name: "prune",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Prune")
			client := kontena.Client{}

			stacks, err := client.StackList()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			for _, stack := range stacks {
				if stack == "core" {
					continue
				}
				if _, err := os.Stat(fmt.Sprintf("./stacks/%s", stack)); os.IsNotExist(err) {
					if err := client.StackRemove(stack); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
			}
			return nil
		},
	}
}
