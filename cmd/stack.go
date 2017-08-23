package cmd

import (
	"os"

	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"

	"github.com/urfave/cli"
)

// StackCommand ...
func StackCommand() cli.Command {
	return cli.Command{
		Name: "stack",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "grid",
				EnvVar: "GRID",
				Usage:  "grid used for installing",
			},
			cli.StringFlag{
				Name:   "filename",
				EnvVar: "KONTENA_FILENAME",
				Usage:  "filename with yaml config",
				Value:  "./kontena.yml",
			},
		},
		Subcommands: []cli.Command{
			stackInstallCommand(),
		},
	}
}

func stackInstallCommand() cli.Command {
	return cli.Command{
		Name: "install",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}

			if err := client.EnsureMasterLogin(); err != nil {
				return cli.NewExitError(err, 1)
			}

			grid := c.Parent().String("grid")
			if client.CurrentGrid().Name == "" || grid != "" {
				if err := client.GridUse(grid); err != nil {
					return cli.NewExitError(err, 1)
				}
			}

			filename := c.Parent().String("filename")
			stack, stackErr := getStack(filename)
			if stackErr != nil {
				return cli.NewExitError(stackErr.Error(), 1)
			}

			if err := client.StackUpgrade(stack); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := client.StackDeploy(stack.Name); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func getStack(file string) (model.KontenaStack, error) {
	var k model.KontenaStack
	if _, err := os.Stat(file); err != nil {
		return k, err
	}
	return model.KontenaLoad(file)
}
