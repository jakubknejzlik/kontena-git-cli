package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
			if client.CurrentGrid().Name == "" || grid != "" {
				if err := client.GridUse(grid); err != nil {
					return err
				}
			}

			if err := installCoreCommand().Run(c); err != nil {
				return err
			}

			if err := installRegistriesCommand().Run(c); err != nil {
				return err
			}

			if err := pruneStacksCommand().Run(c); err != nil {
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
			pruneStacksCommand(),
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

			dc, err := model.KontenaLoad("kontena.yml")
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := client.StackInstallOrUpgrade(dc); err != nil {
				return cli.NewExitError(err, 1)
			}

			utils.Log("deploying stack core")
			if err := client.StackDeploy("core"); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func pruneStacksCommand() cli.Command {
	return cli.Command{
		Name: "prune",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}

			stacks, err := client.StackList()
			if err != nil {
				return err
			}

			for _, stack := range stacks {
				if stack == "core" {
					continue
				}
				if _, err := os.Stat(fmt.Sprintf("./stacks/%s", stack)); os.IsNotExist(err) {
					if err := client.StackRemove(stack); err != nil {
						return err
					}
				}
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
					log.Println(err)
					return err
				}
				if !client.StackExists(stackName) {
					utils.Log("installing stack", stackName)
					dc := getDefaultStack(stackName)
					if err := client.StackInstall(dc); err != nil {
						return cli.NewExitError(err, 1)
					}
				} else {
					if stack, err := getStack(stackName); err == nil {
						utils.Log("upgrading stack", stackName)
						if err := client.StackUpgrade(stack); err != nil {
							return cli.NewExitError(err, 1)
						}
					}
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

func getStack(name string) (model.Kontena, error) {
	var k model.Kontena
	stackConfigPath := fmt.Sprintf("./stacks/%s/kontena.yml", name)
	if _, err := os.Stat(stackConfigPath); err != nil {
		return k, err
	}
	return model.KontenaLoad(stackConfigPath)
}

func getDefaultStack(name string) model.Kontena {
	if stack, err := getStack(name); err == nil {
		return stack
	}
	return model.Kontena{
		Stack:   name,
		Version: "0.0.1",
		Services: map[string]model.KontenaService{
			"web": model.KontenaService{
				Image: "ksdn117/test-page",
				Links: []string{
					"core/internet_lb",
				},
				Secrets: []model.KontenaSecret{
					model.KontenaSecret{
						Secret: "VIRTUAL_HOSTS",
						Name:   "KONTENA_LB_VIRTUAL_HOSTS",
						Type:   "env",
					},
				},
			},
		},
	}
}
