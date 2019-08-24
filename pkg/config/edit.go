package config

import (
	"os/exec"

	"github.com/urfave/cli"
)

func HandleEdit(c *cli.Context) error {
	cmd := exec.Command("open", configFilePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
