package cmd

import (
	"fmt"
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
			stackRemoveCommand(),
		},
	}
}

func stackInstallCommand() cli.Command {
	return cli.Command{
		Name: "install",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "force",
				EnvVar: "KONTENA_FORCE_INSTALL",
				Usage:  "Force stack installation if it doesn't exists in grid",
			},
		},
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

			if c.Bool("force") {
				if !client.StackExists(stack.Name) {
					if err := client.StackInstall(stack); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
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

func stackRemoveCommand() cli.Command {
	return cli.Command{
		Name: "rm",
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

			stack := c.Args().First()

			if stack == "" {
				return cli.NewExitError(fmt.Errorf("provide stack attribute name"), 1)
			}

			if err := client.StackRemove(stack); err != nil {
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
