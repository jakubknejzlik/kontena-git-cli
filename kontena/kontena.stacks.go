package kontena

import (
	"fmt"
	"os"

	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// StackExists ...
func (c *Client) StackExists(stack string) bool {
	res, err := utils.Run("kontena stack ls | awk 'FNR>1{printf \"%s \",$2}'")

	if err != nil {
		return false
	}

	stacks := utils.SplitString(string(res), " ")

	for _, _stack := range stacks {
		if _stack == stack {
			return true
		}
	}

	return false
}

// StackInstallOrUpgrade ...
func (c *Client) StackInstallOrUpgrade(name string, stack model.Compose) error {
	if c.StackExists(name) {
		return c.StackUpgrade(name, stack)
	}
	return c.StackInstall(name, stack)
}

// StackDeploy ...
func (c *Client) StackDeploy(name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack deploy %s", name))
}

// StackInstall ...
func (c *Client) StackInstall(name string, stack model.Compose) error {
	return c.stackAction("install", name, stack)
}

// StackUpgrade ...
func (c *Client) StackUpgrade(name string, stack model.Compose) error {
	return c.stackAction("upgrade", name, stack)
}

func (c *Client) stackAction(action, name string, stack model.Compose) error {
	ds, err := getProcessedCompose(name, stack)
	if err != nil {
		return err
	}
	dsPath, err := ds.ExportTemporary()
	if err != nil {
		return err
	}

	defer os.Remove(dsPath)

	cmd := fmt.Sprintf("kontena stack upgrade %s %s", name, dsPath)
	if action == "install" {
		cmd = fmt.Sprintf("kontena stack install --name %s %s", name, dsPath)
	}
	return utils.RunInteractive(cmd)
}

func getProcessedCompose(name string, s model.Compose) (model.Compose, error) {
	// remap secrets to stack scope (by adding {stack}_)
	for i, service := range s.Services {
		newSecrets := []model.ComposeSecret{}
		for _, secret := range service.Secrets {
			newSecrets = append(newSecrets, model.ComposeSecret{
				Secret: name + "_" + secret.Secret,
				Name:   secret.Name,
				Type:   secret.Type,
			})
		}
		service.Secrets = newSecrets
		s.Services[i] = service
	}

	return s, nil
}
