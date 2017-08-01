package cmd

import (
	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"

	"github.com/urfave/cli"
)

func installRegistriesCommand() cli.Command {
	return cli.Command{
		Name: "registries",
		Action: func(c *cli.Context) error {
			utils.LogSection("Registries")
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
