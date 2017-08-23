package cmd

import (
	"fmt"
	"io/ioutil"
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
		Subcommands: []cli.Command{
			gridInstallCommand(),
		},
	}
}

func gridInstallCommand() cli.Command {
	return cli.Command{
		Name:      "install",
		ArgsUsage: "GRID",
		Action: func(c *cli.Context) error {
			client := kontena.Client{}
			grid := c.Args().First()

			if err := client.EnsureMasterLogin(); err != nil {
				return cli.NewExitError(err, 1)
			}

			if grid == "" {
				return cli.NewExitError("GRID argument not specified", 1)
			}

			if err := client.GridUse(grid); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := installCertificatesCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := installCoreCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := installRegistriesCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := pruneStacksCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := installStacksCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
		// Subcommands: []cli.Command{
		// 	installCoreCommand(),
		// 	installRegistriesCommand(),
		// 	pruneStacksCommand(),
		// 	installStacksCommand(),
		// 	installCertificatesCommand(),
		// },
	}
}

func installStacksCommand() cli.Command {
	return cli.Command{
		Name: "stacks",
		Action: func(c *cli.Context) error {
			utils.LogSection("Stacks")
			client := kontena.Client{}

			stacks, _ := ioutil.ReadDir("./stacks")
			for _, stack := range stacks {
				stackName := stack.Name()
				if err := client.SecretsImport(stackName, fmt.Sprintf("./stacks/%s/secrets.yml", stackName)); err != nil {
					return cli.NewExitError(err, 1)
				}
				if !client.StackExists(stackName) {
					utils.Log("installing stack", stackName)
					dc := getDefaultStack(stackName, client.SecretExists("VIRTUAL_HOSTS", stackName))
					if err := client.StackInstall(dc); err != nil {
						return cli.NewExitError(err, 1)
					}
				} else {
					if stack, err := getStackFromGrid(stackName); err == nil {
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

func getStackFromGrid(name string) (model.KontenaStack, error) {
	var k model.KontenaStack
	stackConfigPath := fmt.Sprintf("./stacks/%s/kontena.yml", name)
	if _, err := os.Stat(stackConfigPath); err != nil {
		return k, err
	}
	return model.KontenaLoad(stackConfigPath)
}

func getDefaultStack(name string, hasHost bool) model.KontenaStack {
	secrets := []model.KontenaSecret{}
	links := []string{}

	if hasHost {
		hostSecret := model.KontenaSecret{
			Secret: "VIRTUAL_HOSTS",
			Name:   "KONTENA_LB_VIRTUAL_HOSTS",
			Type:   "env",
		}
		secrets = append(secrets, hostSecret)
		links = append(links, "core/internet_lb")
	}

	return model.KontenaStack{
		Name:    name,
		Version: "0.0.1",
		Services: map[string]model.KontenaService{
			"web": model.KontenaService{
				Image:   "ksdn117/test-page",
				Links:   links,
				Secrets: secrets,
			},
		},
	}
}
