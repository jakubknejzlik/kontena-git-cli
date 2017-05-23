package kontena

import (
	"fmt"

	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// ServiceDeploy ...
func (c *Client) ServiceDeploy(stack, service string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena service deploy %s/%s", stack, service))
}

// ServiceExec ...
func (c *Client) ServiceExec(stack, service, command string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena service exec %s/%s %s", stack, service, command))
}
