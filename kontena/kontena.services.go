package kontena

import (
	"fmt"
	"strings"

	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// ServiceCreate ...
func (c *Client) ServiceCreate(name string, service model.KontenaService) error {
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
	cmd = append(cmd, name)
	cmd = append(cmd, service.Image)

	utils.Log("creating service", name)
	return utils.RunInteractive(strings.Join(cmd, " "))
}

// ServiceDeploy ...
func (c *Client) ServiceDeploy(service string) error {
	utils.Log("deploying service", service)
	return utils.RunInteractive(fmt.Sprintf("kontena service deploy %s", service))
}

// ServiceInStackDeploy ...
func (c *Client) ServiceInStackDeploy(stack, service string) error {
	return c.ServiceDeploy(stack + "/" + service)
}

// ServiceExec ...
func (c *Client) ServiceExec(service, command string) ([]byte, error) {
	return utils.Run(fmt.Sprintf("kontena service exec %s %s", service, command))
}

// ServiceInStackExec ...
func (c *Client) ServiceInStackExec(stack, service, command string) ([]byte, error) {
	return c.ServiceExec(stack+"/"+service, command)
}

// ServiceRemove ...
func (c *Client) ServiceRemove(service string) error {
	utils.Log("removing service", service)
	return utils.RunInteractive(fmt.Sprintf("kontena service rm --force %s", service))
}
