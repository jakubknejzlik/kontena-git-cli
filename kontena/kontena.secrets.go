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

	oldSecrets, err := c.getSecrets()
	if err != nil {
		return err
	}
	for _, secret := range oldSecrets {
		if strings.HasPrefix(secret, stack+"_") {
			utils.Log("removing secret", strings.Replace(secret, stack+"_", stack+":", 1))
			c.removeSecret(secret)
		}
	}

	for key, value := range secrets {
		secretKey := fmt.Sprintf("%s_%s", stack, key)
		utils.Log("adding secret", stack+":"+key)
		cmd := fmt.Sprintf("kontena vault write %s %s", secretKey, value)
		if err := utils.RunInteractive(cmd); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) removeSecret(secret string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena vault rm --force %s", secret))
}

func (c *Client) getSecrets() ([]string, error) {
	data, err := utils.Run("kontena vault ls -q")
	return utils.SplitString(string(data), "\n"), err
}

func (c *Client) getSecret(name string) (string, error) {
	value, err := utils.Run(fmt.Sprintf("kontena vault read --value %s", name))
	return string(value), err
}
