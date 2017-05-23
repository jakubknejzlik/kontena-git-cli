package kontena

import (
	"fmt"

	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// RegistryExists ...
func (c *Client) RegistryExists(name string) bool {
	registries, err := c.CurrentRegistries()
	if err != nil {
		return false
	}

	for _, registry := range registries {
		if registry == name {
			return true
		}
	}

	return false
}

// RegistryAdd ...
func (c *Client) RegistryAdd(registry model.Registry) error {
	cmd := fmt.Sprintf("kontena external-registry add --username %s --password %s --email %s https://%s/v2/", registry.User, registry.Password, registry.Email, registry.Name)
	return utils.RunInteractive(cmd)
}

// RegistryRemove ...
func (c *Client) RegistryRemove(name string) error {
	cmd := fmt.Sprintf("kontena external-registry rm --force %s", name)
	return utils.RunInteractive(cmd)
}

// CurrentRegistries ...
func (c *Client) CurrentRegistries() ([]string, error) {
	res, err := utils.Run("kontena external-registry ls | awk 'FNR>1{printf \"%s \",$1}'")
	if err != nil {
		return nil, err
	}

	return utils.SplitString(string(res), " "), nil
}
