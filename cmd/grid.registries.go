package cmd

import (
	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"

	"github.com/urfave/cli"
)

func installRegistriesCommand() cli.Command {
	return cli.Command{
		Name: "registries",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Registries")
			client := kontena.Client{}

			currentRegistries, listErr := client.RegistryList()
			if listErr != nil {
				return cli.NewExitError(listErr, 1)
			}

			registries, loadErr := model.RegistriesLoad("registries.yml")
			if loadErr != nil {
				return cli.NewExitError(loadErr, 1)
			}
			registryNames := []string{}
			for _, reg := range registries {
				registryNames = append(registryNames, reg.Name)
			}

			for _, regName := range currentRegistries {
				if utils.ArrayOfStringsContains(currentRegistries, regName) && !utils.ArrayOfStringsContains(registryNames, regName) {
					if err := client.RegistryRemove(regName); err != nil {
						return cli.NewExitError(err.Error(), 1)
					}
				}
			}

			for _, registry := range registries {
				if !utils.ArrayOfStringsContains(currentRegistries, registry.Name) && utils.ArrayOfStringsContains(registryNames, registry.Name) {
					if err := client.RegistryAdd(registry); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
			}

			return nil
		},
	}
}
