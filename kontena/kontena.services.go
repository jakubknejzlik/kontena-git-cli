package kontena

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// ServiceCreate ...
func (c *Client) ServiceCreate(name string, service model.KontenaService) error {
	return c.ServiceCreateInGrid(c.CurrentGrid().Name, name, service)
}

// ServiceCreateInGrid ...
func (c *Client) ServiceCreateInGrid(grid, name string, service model.KontenaService) error {
	cmd := []string{`kontena service create`}
	if service.Instances != nil && *service.Instances > 0 {
		cmd = append(cmd, `--instances `+string(*service.Instances))
	}
	if service.Command != "" {
		cmd = append(cmd, `--cmd `+service.Command)
	}
	for _, value := range service.Environment {
		cmd = append(cmd, `-e "`+value+`"`)
	}
	for _, value := range service.Links {
		cmd = append(cmd, `-l "`+value+`"`)
	}
	for _, value := range service.Volumes {
		cmd = append(cmd, `-v "`+value+`"`)
	}
	for _, value := range service.Ports {
		cmd = append(cmd, `-p "`+value+`"`)
	}
	if service.Deploy.Strategy != "" {
		cmd = append(cmd, `--deploy `+service.Deploy.Strategy)
	}
	cmd = append(cmd, `--grid `+grid)
	cmd = append(cmd, name)
	cmd = append(cmd, service.Image)

	utils.Log("creating service", name, "in grid", grid)
	return utils.RunInteractive(strings.Join(cmd, " "))
}

// ServiceDeploy ...
func (c *Client) ServiceDeploy(service string) error {
	utils.Log("deploying service", service)
	return utils.RunInteractive(fmt.Sprintf("kontena service deploy %s", service))
}

// ServiceDeployInGrid ...
func (c *Client) ServiceDeployInGrid(grid, service string) error {
	utils.Log("deploying service", service, "in grid", grid)
	return utils.RunInteractive(fmt.Sprintf("kontena service deploy --grid %s %s", grid, service))
}

// ServiceInStackDeploy ...
func (c *Client) ServiceInStackDeploy(stack, service string) error {
	return c.ServiceDeploy(stack + "/" + service)
}

// ServiceInStackInGridDeploy ...
func (c *Client) ServiceInStackInGridDeploy(grid, stack, service string) error {
	return c.ServiceDeployInGrid(grid, stack+"/"+service)
}

// ServiceExec ...
func (c *Client) ServiceExec(service, command string) ([]byte, error) {
	return utils.Run(fmt.Sprintf("kontena service exec %s %s", service, command))
}

// ServiceExecInGrid ...
func (c *Client) ServiceExecInGrid(grid, service, command string) ([]byte, error) {
	return utils.Run(fmt.Sprintf("kontena service exec --grid %s %s %s", grid, service, command))
}

// ServiceExecCommand ...
func (c *Client) ServiceExecCommand(service, command string) *exec.Cmd {
	return utils.RunCommand(fmt.Sprintf("kontena service exec %s %s", service, command))
}

// ServiceExecInGridCommand ...
func (c *Client) ServiceExecInGridCommand(grid, service, command string) *exec.Cmd {
	return utils.RunCommand(fmt.Sprintf("kontena service exec --grid %s %s %s", grid, service, command))
}

// ServiceInStackExec ...
func (c *Client) ServiceInStackExec(stack, service, command string) ([]byte, error) {
	return c.ServiceExec(stack+"/"+service, command)
}

// ServiceInStackInGridExec ...
func (c *Client) ServiceInStackInGridExec(grid, stack, service, command string) ([]byte, error) {
	return c.ServiceExecInGrid(grid, stack+"/"+service, command)
}

// ServiceInStackExecCommand ...
func (c *Client) ServiceInStackExecCommand(stack, service, command string) *exec.Cmd {
	return c.ServiceExecCommand(stack+"/"+service, command)
}

// ServiceInStackInGridExecCommand ...
func (c *Client) ServiceInStackInGridExecCommand(grid, stack, service, command string) *exec.Cmd {
	return c.ServiceExecInGridCommand(grid, stack+"/"+service, command)
}

// ServiceRemove ...
func (c *Client) ServiceRemove(service string) error {
	utils.Log("removing service", service)
	return utils.RunInteractive(fmt.Sprintf("kontena service rm --force %s", service))
}

// ServiceList ...
func (c *Client) ServiceList() ([]string, error) {
	data, err := utils.Run("kontena service ls -q")
	if err != nil {
		return []string{}, err
	}
	return utils.SplitString(string(data), "\n"), nil
}

// ServiceListInGrid ...
func (c *Client) ServiceListInGrid(grid string) ([]string, error) {
	data, err := utils.Run(fmt.Sprintf("kontena service ls --grid %s -q", grid))
	if err != nil {
		return []string{}, err
	}
	return utils.SplitString(string(data), "\n"), nil
}

// ServiceExists ...
func (c *Client) ServiceExists(stack, service string) (bool, error) {
	grid := c.CurrentGrid().Name
	return c.ServiceExistsInGrid(grid, stack, service)
}

// ServiceExistsInGrid ...
func (c *Client) ServiceExistsInGrid(grid, stack, service string) (bool, error) {
	services, err := c.ServiceListInGrid(grid)
	if err != nil {
		return false, err
	}
	if stack == "" {
		stack = "null"
	}
	return utils.ArrayOfStringsContains(services, grid+"/"+stack+"/"+service), nil
}

// ServiceLogs ...
func (c *Client) ServiceLogs(service string) (string, error) {
	data, err := utils.Run(fmt.Sprintf("kontena service logs %s", service))
	return string(data), err
}

// ServiceInStackLogs ...
func (c *Client) ServiceInStackLogs(stack, service string) (string, error) {
	data, err := utils.Run(fmt.Sprintf("kontena service logs %s/%s", stack, service))
	return string(data), err
}

// ServiceInStackInGridLogs ...
func (c *Client) ServiceInStackInGridLogs(grid, stack, service string) (string, error) {
	data, err := utils.Run(fmt.Sprintf("kontena service logs --grid %s %s/%s", grid, stack, service))
	return string(data), err
}
