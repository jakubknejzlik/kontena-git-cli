package kontena

import (
	"fmt"
	"strings"

	"github.com/inloop/goclitools"
	"github.com/urfave/cli"
)

// Grid ...
type Grid struct {
	Name string
}

// GridUse ...
func (c *Client) GridUse(grid string) error {
	if grid == "" {
		return cli.NewExitError("grid must be specified", 1)
	}
	return goclitools.RunInteractive(fmt.Sprintf("kontena grid use %s", grid))
}

// CurrentGrid ...
func (c *Client) CurrentGrid() Grid {
	var g Grid
	res, err := goclitools.Run("kontena grid current --name")
	if err != nil {
		return g
	}
	return Grid{Name: strings.Trim(string(res), "\n")}
}
