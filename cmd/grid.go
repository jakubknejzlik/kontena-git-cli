package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"

	"github.com/urfave/cli"
)

// GridCommand ...
func GridCommand() cli.Command {
	return cli.Command{
		Name: "grid",
		Subcommands: []cli.Command{
			gridInstallCommand(),
			gridDeployCommand(),
			gridCleanupCommand(),
		},
	}
}

func gridInstallCommand() cli.Command {
	return cli.Command{
		Name:      "install",
		ArgsUsage: "GRID",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "deploy",
				Usage: "automatically deploy all services",
			},
		},
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

			// install stack if not already installed to be able to run installCertificatesCommand
			if client.StackExists("core") == false {
				if err := installCoreCommand().Run(c); err != nil {
					return cli.NewExitError(err, 1)
				}
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

			if err := installOrUpgradeStacksCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			if c.Bool("deploy") {
				if err := deployStacksCommand().Run(c); err != nil {
					return cli.NewExitError(err, 1)
				}
			}

			return nil
		},
	}
}

func gridDeployCommand() cli.Command {
	return cli.Command{
		Name:        "deploy",
		ArgsUsage:   "GRID",
		Description: "Deploy all stacks in grid",
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

			if err := deployStacksCommand().Run(c); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func gridCleanupCommand() cli.Command {
	return cli.Command{
		Name:        "cleanup",
		ArgsUsage:   "GRID",
		Description: "Cleanup grid (renew certificates, etc.)",
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

			return nil
		},
	}
}

func installOrUpgradeStacksCommand() cli.Command {
	return cli.Command{
		Name: "stacks",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Installing/upgrading stacks")
			client := kontena.Client{}

			goclitools.Log("Updating secrets...")

			currentSecrets, err := client.SecretList()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			secretsImports := map[string]string{}

			stacks, _ := ioutil.ReadDir("./stacks")
			for _, stack := range stacks {
				stackName := stack.Name()
				secretsImport, err := client.CreateSecretsImport(stackName, fmt.Sprintf("./stacks/%s/secrets.yml", stackName), currentSecrets)
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				for key, value := range secretsImport {
					secretsImports[key] = value
				}
			}

			client.SecretsImportInGrid(secretsImports)

			goclitools.Log("Secrets updated")

			currentStacks, err := client.StackList()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			currentStacksMap := map[string]bool{}
			for _, stackName := range currentStacks {
				currentStacksMap[stackName] = true
			}

			for _, stack := range stacks {
				stackName := stack.Name()

				if !currentStacksMap[stackName] {
					goclitools.Log("installing stack", stackName)
					dc, stackErr := getStackFromGrid(stackName)
					if stackErr != nil {
						dc = getDefaultStack(stackName)
					}

					if err := client.StackInstall(dc); err != nil {
						return cli.NewExitError(err, 1)
					}
				} else {
					if stack, err := getStackFromGrid(stackName); err == nil {
						goclitools.Log("upgrading stack", stackName)
						if err := client.StackUpgrade(stack); err != nil {
							return cli.NewExitError(err, 1)
						}
					}
				}
			}

			return nil
		},
	}
}

func deployStacksCommand() cli.Command {
	return cli.Command{
		Name: "stacks",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Deploying stacks")
			client := kontena.Client{}

			stacks, _ := ioutil.ReadDir("./stacks")
			for _, stack := range stacks {
				stackName := stack.Name()

				goclitools.Log("deploying stack", stackName)
				if err := client.StackDeploy(stackName); err != nil {
					return cli.NewExitError(err, 1)
				}
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

func getDefaultStack(name string) model.KontenaStack {
	return model.KontenaStack{
		Name:    name,
		Version: "0.0.1",
		Services: map[string]model.KontenaService{
			"dummy": model.KontenaService{Image: "tutum/hello-world:latest"},
		},
	}
}
