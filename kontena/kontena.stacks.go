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

// StackListInGrid ...
func (c *Client) StackListInGrid(grid string) ([]string, error) {
	var list []string
	res, err := utils.Run(fmt.Sprintf("kontena stack ls --grid %s -q", grid))

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

// StackExistsInGrid ...
func (c *Client) StackExistsInGrid(grid, stack string) bool {
	stacks, err := c.StackListInGrid(grid)
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

// StackInstallOrUpgradeInGrid ...
func (c *Client) StackInstallOrUpgradeInGrid(grid string, stack model.KontenaStack) error {
	if c.StackExistsInGrid(grid, stack.Name) {
		return c.StackUpgradeInGrid(grid, stack)
	}
	return c.StackInstallInGrid(grid, stack)
}

// StackDeploy ...
func (c *Client) StackDeploy(name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack deploy %s", name))
}

// StackDeployInGrid ...
func (c *Client) StackDeployInGrid(grid, name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack deploy --grid %s %s", grid, name))
}

// StackInstall ...
func (c *Client) StackInstall(stack model.KontenaStack) error {
	return c.stackAction("install", stack.Name, stack)
}

// StackInstallInGrid ...
func (c *Client) StackInstallInGrid(grid string, stack model.KontenaStack) error {
	return c.stackActionInGrid("install", grid, stack.Name, stack)
}

// StackUpgrade ...
func (c *Client) StackUpgrade(stack model.KontenaStack) error {
	return c.stackAction("upgrade", stack.Name, stack)
}

// StackUpgradeInGrid ...
func (c *Client) StackUpgradeInGrid(grid string, stack model.KontenaStack) error {
	return c.stackActionInGrid("upgrade", grid, stack.Name, stack)
}

// StackRemove ...
func (c *Client) StackRemove(name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack remove --force %s", name))
}

// StackRemoveFromGrid ...
func (c *Client) StackRemoveFromGrid(grid, name string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena stack remove --grid %s --force %s", grid, name))
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

func (c *Client) stackActionInGrid(grid, action, name string, stack model.KontenaStack) error {
	dsPath, err := stack.ExportTemporary(true)
	if err != nil {
		return err
	}

	defer os.Remove(dsPath)

	cmd := fmt.Sprintf("kontena stack upgrade --force --grid %s --no-deploy %s %s", grid, name, dsPath)
	if action == "install" {
		cmd = fmt.Sprintf("kontena stack install --grid %s --name %s %s", grid, name, dsPath)
	}
	return utils.RunInteractive(cmd)
}
