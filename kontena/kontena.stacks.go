package kontena

import (
	"fmt"
	"os"

	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// StackList ...
func (c *Client) StackList() ([]string, error) {
	var list []string
	res, err := utils.Run("kontena stack ls -q")

	if err != nil {
		return list, err
	}

	return utils.SplitString(string(res), "\n"), nil
}

// StackExists ...
func (c *Client) StackExists(stack string) bool {
	stacks, err := c.StackList()
	if err != nil {
		return false
	}

	for _, _stack := range stacks {
		if _stack == stack {
			return true
		}
	}

	return false
}

// StackInstallOrUpgrade ...
func (c *Client) StackInstallOrUpgrade(stack model.KontenaStack) error {
	if c.StackExists(stack.Name) {
		return c.StackUpgrade(stack)
	}
	return c.StackInstall(stack)
}

// StackDeploy ...
func (c *Client) StackDeploy(name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack deploy %s", name))
}

// StackInstall ...
func (c *Client) StackInstall(stack model.KontenaStack) error {
	return c.stackAction("install", stack.Name, stack)
}

// StackUpgrade ...
func (c *Client) StackUpgrade(stack model.KontenaStack) error {
	return c.stackAction("upgrade", stack.Name, stack)
}

// StackRemove ...
func (c *Client) StackRemove(name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack remove --force %s", name))
}

func (c *Client) stackAction(action, name string, stack model.KontenaStack) error {
	dsPath, err := stack.ExportTemporary(true)
	if err != nil {
		return err
	}

	defer os.Remove(dsPath)

	cmd := fmt.Sprintf("kontena stack upgrade --no-deploy %s %s", name, dsPath)
	if action == "install" {
		cmd = fmt.Sprintf("kontena stack install --name %s %s", name, dsPath)
	}
	return utils.RunInteractive(cmd)
}
