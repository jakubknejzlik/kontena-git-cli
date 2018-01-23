package utils

import (
	"os/exec"
)

// RunCommand ...
func RunCommand(cmd string) *exec.Cmd {
	return exec.Command("sh", "-c", cmd)
}
