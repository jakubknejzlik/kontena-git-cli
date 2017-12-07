package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/jakubknejzlik/kontena-git-cli/kontena"

	"github.com/urfave/cli"
)

// ServiceCommand ...
func ServiceCommand() cli.Command {
	return cli.Command{
		Name: "service",
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
			stackRunCommand(),
		},
	}
}

func stackRunCommand() cli.Command {
	return cli.Command{
		Name: "exec",
		Action: func(c *cli.Context) error {
			grid := c.String("grid")

			args := c.Args()
			cmd := strings.Join(args.Tail(), " ")

			serviceName := strings.Split(args.First(), "/")

			return runCommandInStack(grid, serviceName[0], serviceName[1], cmd)
		},
	}
}

func runCommandInStack(grid, stack, service, command string) error {

	client := kontena.Client{}

	if err := client.EnsureMasterLogin(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	if grid != "" {
		cmd = client.ServiceInStackInGridExecCommand(grid, stack, service, command)
	} else {
		cmd = client.ServiceInStackExecCommand(stack, service, command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
