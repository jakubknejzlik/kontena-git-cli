package utils

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

// Run ...
func Run(cmd string) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	data, err := command.Output()
	if err != nil {
		return data, cli.NewExitError(err.Error(), 1)
	}
	return data, nil
}

// RunWithInput ...
func RunWithInput(cmd string, input []byte) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	command.Stdin = bytes.NewReader(input)
	data, err := command.Output()
	if err != nil {
		return data, cli.NewExitError(err.Error(), 1)
	}
	return data, nil
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
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}
