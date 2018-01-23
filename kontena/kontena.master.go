package kontena

import (
	"fmt"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

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
	res, err := goclitools.Run("kontena master current")
	if err != nil {
		return m
	}

	masterPair := utils.SplitString(string(res), " ")
	m = Master{masterPair[0], masterPair[1]}

	return m
}

// EnsureMasterLogin ...
func (c *Client) EnsureMasterLogin() error {
	masterURL := utils.Getenv("KONTENA_MASTER_URL", "")

	currentMasterURL := c.CurrentMaster().url

	if masterURL != "" && currentMasterURL != masterURL || currentMasterURL == "" {
		masterURL = utils.GetenvStrict("KONTENA_MASTER_URL")
		token := utils.GetenvStrict("KONTENA_TOKEN")
		fmt.Println("logging to master", masterURL)
		return c.MasterLogin(masterURL, token)
	}

	// fmt.Println("logged to master", c.CurrentMaster().url)
	return nil
}

// MasterLogin ...
func (c *Client) MasterLogin(masterURL, token string) error {
	return goclitools.RunInteractive(fmt.Sprintf("kontena master login --skip-grid-auto-select --token %s %s", token, masterURL))
}
