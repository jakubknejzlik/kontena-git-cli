package utils

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

// Run ...
func Run(cmd string) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	return command.Output()
}

// RunInteractive ...
func RunInteractive(cmd string) error {
	if os.Getenv("DEBUG") != "" {
		log.Println(cmd)
	}
	command := exec.Command("sh", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	// command.Dir = cwd
	err := command.Run()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}
