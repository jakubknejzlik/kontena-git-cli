package cmd

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"

	"github.com/urfave/cli"
)

// GridCommand ...
func GridCommand() cli.Command {
	return cli.Command{
		Name: "grid",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "grid",
				EnvVar: "GRID",
				Usage:  "grid used for installing",
			},
		},
		Subcommands: []cli.Command{
			installCommand(),
		},
	}
}

func installCommand() cli.Command {
	return cli.Command{
		Name: "install",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}

			if err := client.EnsureMasterLogin(); err != nil {
				return err
			}

			grid := c.Parent().String("grid")
			if err := client.GridUse(grid); err != nil {
				return err
			}

			if err := installCoreCommand().Run(c); err != nil {
				return err
			}

			if err := installRegistriesCommand().Run(c); err != nil {
				return err
			}

			if err := installStacksCommand().Run(c); err != nil {
				return err
			}

			return nil
		},
		Subcommands: []cli.Command{
			installCoreCommand(),
			installRegistriesCommand(),
			installStacksCommand(),
		},
	}
}

func installRegistriesCommand() cli.Command {
	return cli.Command{
		Name: "registries",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}

			currentRegistries, err := client.CurrentRegistries()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			for _, regName := range currentRegistries {
				if client.RegistryExists(regName) {
					client.RegistryRemove(regName)
				}
			}

			registries, err := model.RegistriesLoad("registries.yml")
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			for _, registry := range registries {
				if !client.RegistryExists(registry.Name) {
					if err := client.RegistryAdd(registry); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
			}

			return nil
		},
	}
}

func installCoreCommand() cli.Command {
	return cli.Command{
		Name: "core",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}

			dc, err := model.ComposeLoad("kontena.yml")
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := client.StackInstallOrUpgrade("core", dc); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func installStacksCommand() cli.Command {
	return cli.Command{
		Name: "stacks",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}

			stacks, _ := ioutil.ReadDir("./stacks")
			for _, stack := range stacks {
				stackName := stack.Name()
				if err := client.SecretsImport(stackName, fmt.Sprintf("./stacks/%s/secrets.yml", stackName)); err != nil {
					return err
				}
				if !client.StackExists(stackName) {
					utils.Log("installing stack", stackName)
					dc := defaultStack(stackName)
					if err := client.StackInstallOrUpgrade(stackName, dc); err != nil {
						return cli.NewExitError(err, 1)
					}
				} else {
					utils.Log("deploying stack", stackName)
					if err := client.StackDeploy(stackName); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
				time.Sleep(time.Second * 3)
			}
			return nil
		},
	}
}

func defaultStack(name string) model.Compose {
	return model.Compose{
		Stack:   name,
		Version: "0.0.1",
		Services: map[string]model.ComposeService{
			"web": model.ComposeService{
				Image: "ksdn117/test-page",
				Links: []string{
					"core/internet_lb",
				},
				Secrets: []model.ComposeSecret{
					model.ComposeSecret{
						Secret: "VIRTUAL_HOSTS",
						Name:   "KONTENA_LB_VIRTUAL_HOSTS",
						Type:   "env",
					},
				},
			},
		},
	}
}
