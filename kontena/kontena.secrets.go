package kontena

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// CreateSecretsImport ...
func (c *Client) CreateSecretsImport(stack, path string, currentSecrets []model.Secret) (map[string]string, error) {
	result := map[string]string{}

	var secrets map[string]string
	data, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return result, err
	}

	yaml.Unmarshal(data, &secrets)

	secretsToRemove := []string{}

	for _, secret := range currentSecrets {
		secretKey := strings.Replace(secret.Name, stack+"_", "", 1)
		if strings.HasPrefix(secret.Name, stack+"_") && secrets[secretKey] == "" {
			secretsToRemove = append(secretsToRemove, secret.Name)
		}
	}

	for _, key := range secretsToRemove {
		result[key] = ""
	}
	for key, value := range secrets {
		result[stack+"_"+key] = value
	}

	return result, nil
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
	_, err := goclitools.RunWithInput(cmd, []byte(value))
	return err
}

// SecretWriteToGrid ...
func (c *Client) SecretWriteToGrid(grid, secret, value string) error {
	cmd := fmt.Sprintf("kontena vault update --grid %s -u %s", grid, secret)
	_, err := goclitools.RunWithInput(cmd, []byte(value))
	return err
}

// SecretRemove ...
func (c *Client) SecretRemove(secret string) error {
	return goclitools.RunInteractive(fmt.Sprintf("kontena vault rm --force %s", secret))
}

// SecretRemoveFromGrid ...
func (c *Client) SecretRemoveFromGrid(grid, secret string) error {
	return goclitools.RunInteractive(fmt.Sprintf("kontena vault rm --grid %s --force %s", grid, secret))
}

// SecretList ...
func (c *Client) SecretList() ([]model.Secret, error) {
	data, err := goclitools.Run("kontena vault ls -l")
	if err != nil {
		return []model.Secret{}, err
	}
	rows := utils.SplitString(string(data), "\n")
	return model.SecretParseList(rows)
}

// SecretListInGrid ...
func (c *Client) SecretListInGrid(grid string) ([]model.Secret, error) {
	data, err := goclitools.Run(fmt.Sprintf("kontena vault ls -l --grid %s", grid))
	if err != nil {
		return []model.Secret{}, err
	}
	rows := utils.SplitString(string(data), "\n")
	return model.SecretParseList(rows)
}

// SecretValue ...
func (c *Client) SecretValue(name string) (string, error) {
	value, err := goclitools.Run(fmt.Sprintf("kontena vault read --value %s", name))
	return string(value), err
}

// SecretValueInGrid ...
func (c *Client) SecretValueInGrid(grid, name string) (string, error) {
	value, err := goclitools.Run(fmt.Sprintf("kontena vault read --grid %s --value %s", grid, name))
	return string(value), err
}

// SecretsImportInGrid ...
func (c *Client) SecretsImportInGrid(secrets map[string]string) error {
	secretsInput, err := yaml.Marshal(secrets)
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("kontena vault import --force --empty-is-null")
	_, err = goclitools.RunWithInput(cmd, secretsInput)
	return err
}
