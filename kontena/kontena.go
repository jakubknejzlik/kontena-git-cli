package kontena

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// Client ...
type Client struct {
	currentMaster *Master
}

// Master ...
type Master struct {
	name string
	url  string
}

// CurrentMaster ...
func (c *Client) CurrentMaster() Master {

	if c.currentMaster != nil {
		return *c.currentMaster
	}

	var m Master
	res, err := utils.Run("kontena master current")
	if err != nil {
		return m
	}

	masterPair := utils.SplitString(string(res), " ")
	m = Master{masterPair[0], masterPair[1]}

	return m
}

// EnsureMasterLogin ...
func (c *Client) EnsureMasterLogin() error {
	masterURL := utils.GetenvStrict("KONTENA_MASTER_URL")
	token := utils.GetenvStrict("KONTENA_TOKEN")

	if c.CurrentMaster().url != masterURL {
		fmt.Println("logging to master", masterURL)
		return c.MasterLogin(masterURL, token)
	}

	// fmt.Println("logged to master", c.CurrentMaster().url)
	return nil
}

// MasterLogin ...
func (c *Client) MasterLogin(masterURL, token string) error {
	return utils.RunInteractive(fmt.Sprintf("kontena master login --skip-grid-auto-select --token %s %s", token, masterURL))
}

// GridUse ...
func (c *Client) GridUse(grid string) error {
	if grid == "" {
		return cli.NewExitError("grid must be specified", 1)
	}
	return utils.RunInteractive(fmt.Sprintf("kontena grid use %s", grid))
}
