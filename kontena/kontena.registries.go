package kontena

import (
	"fmt"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// RegistryExists ...
func (c *Client) RegistryExists(name string) (bool, error) {
	registries, err := c.RegistryList()
	if err != nil {
		return false, err
	}

	return utils.ArrayOfStringsContains(registries, name), nil
}

// RegistryExistsInGrid ...
func (c *Client) RegistryExistsInGrid(grid, name string) (bool, error) {
	registries, err := c.RegistryListInGrid(grid)
	if err != nil {
		return false, err
	}

	return utils.ArrayOfStringsContains(registries, name), nil
}

// RegistryAdd ...
func (c *Client) RegistryAdd(registry model.Registry) error {
	cmd := fmt.Sprintf("kontena external-registry add --username %s --password %s --email %s https://%s/v2/", registry.User, registry.Password, registry.Email, registry.Name)
	return goclitools.RunInteractive(cmd)
}

// RegistryAddToGrid ...
func (c *Client) RegistryAddToGrid(grid string, registry model.Registry) error {
	cmd := fmt.Sprintf("kontena external-registry add --grid %s --username %s --password %s --email %s https://%s/v2/", grid, registry.User, registry.Password, registry.Email, registry.Name)
	return goclitools.RunInteractive(cmd)
}

// RegistryRemove ...
func (c *Client) RegistryRemove(name string) error {
	cmd := fmt.Sprintf("kontena external-registry rm --force %s", name)
	return goclitools.RunInteractive(cmd)
}

// RegistryRemoveFromGrid ...
func (c *Client) RegistryRemoveFromGrid(grid, name string) error {
	cmd := fmt.Sprintf("kontena external-registry rm --grid %s --force %s", grid, name)
	return goclitools.RunInteractive(cmd)
}

// RegistryList ...
func (c *Client) RegistryList() ([]string, error) {
	res, err := goclitools.Run("kontena external-registry ls | awk 'FNR>1{printf \"%s \",$1}'")
	if err != nil {
		return nil, err
	}

	return utils.SplitString(string(res), " "), nil
}

// RegistryListInGrid ...
func (c *Client) RegistryListInGrid(grid string) ([]string, error) {
	res, err := goclitools.Run(fmt.Sprintf("kontena external-registry ls --grid %s | awk 'FNR>1{printf \"%%s \",$1}'", grid))
	if err != nil {
		return nil, err
	}

	return utils.SplitString(string(res), " "), nil
}
