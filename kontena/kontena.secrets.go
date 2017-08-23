package kontena

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/jakubknejzlik/kontena-git-cli/utils"
	"github.com/urfave/cli"
)

// SecretsImport ...
func (c *Client) SecretsImport(stack, path string) error {
	var secrets map[string]string

	data, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return cli.NewExitError(err, 1)
	}

	yaml.Unmarshal(data, &secrets)

	oldSecrets, err := c.SecretList()
	if err != nil {
		return err
	}
	for _, secret := range oldSecrets {
		if strings.HasPrefix(secret, stack+"_") {
			utils.Log("removing secret", strings.Replace(secret, stack+"_", stack+":", 1))
			c.SecretRemove(secret)
		}
	}

	for key, value := range secrets {
		secretKey := fmt.Sprintf("%s_%s", stack, key)
		utils.Log("adding secret", stack+":"+key)
		if err := c.SecretWrite(secretKey, value); err != nil {
			return err
		}
	}

	return nil
}

// SecretExists ...
func (c *Client) SecretExists(name, stack string) bool {
	value, _ := c.SecretValue(stack + "_" + name)
	return value != ""
}

// SecretExistsInGrid ...
func (c *Client) SecretExistsInGrid(grid, name, stack string) bool {
	value, _ := c.SecretValueInGrid(grid, stack+"_"+name)
	return value != ""
}

// SecretWrite ...
func (c *Client) SecretWrite(secret, value string) error {
	cmd := fmt.Sprintf("kontena vault update -u %s", secret)
	_, err := utils.RunWithInput(cmd, []byte(value))
	return err
}

// SecretWriteToGrid ...
func (c *Client) SecretWriteToGrid(grid, secret, value string) error {
	cmd := fmt.Sprintf("kontena vault update --grid %s -u %s", grid, secret)
	_, err := utils.RunWithInput(cmd, []byte(value))
	return err
}

// SecretRemove ...
func (c *Client) SecretRemove(secret string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena vault rm --force %s", secret))
}

// SecretRemoveFromGrid ...
func (c *Client) SecretRemoveFromGrid(grid, secret string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena vault rm --grid %s --force %s", grid, secret))
}

// SecretList ...
func (c *Client) SecretList() ([]string, error) {
	data, err := utils.Run("kontena vault ls -q")
	if err != nil {
		return []string{}, err
	}
	return utils.SplitString(string(data), "\n"), nil
}

// SecretListInGrid ...
func (c *Client) SecretListInGrid(grid string) ([]string, error) {
	data, err := utils.Run(fmt.Sprintf("kontena vault ls --grid %s -q", grid))
	if err != nil {
		return []string{}, err
	}
	return utils.SplitString(string(data), "\n"), nil
}

// SecretValue ...
func (c *Client) SecretValue(name string) (string, error) {
	value, err := utils.Run(fmt.Sprintf("kontena vault read --value %s", name))
	return string(value), err
}

// SecretValueInGrid ...
func (c *Client) SecretValueInGrid(grid, name string) (string, error) {
	value, err := utils.Run(fmt.Sprintf("kontena vault read --grid %s --value %s", grid, name))
	return string(value), err
}
