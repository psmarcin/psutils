package config

import (
	"psutils/pkg/utils"

	"github.com/urfave/cli"
)

var runner utils.Runner

func init() {
	runner = utils.RealRunner{}
}

func HandleEdit(c *cli.Context) error {
	_, err := runner.Run("open", configFilePath)
	if err != nil {
		return err
	}
	return err
}
