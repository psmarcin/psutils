package accounting

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
	"os"
	"path"
	"psutils/pkg/config"
	"time"
)

type reformatConfig struct {
	filePath string
	date     time.Time
	name     string
	types    string
}

var l log.Logger

func init() {
	l = log.With("package", "invoice")
}

// Handler holds all logic about create subcommand
func Handler(c *cli.Context) error {
	conf := config.Load()
	date, err := time.Parse(conf.Other.MontDateFormat, c.String("date"))
	if err != nil {
		return err
	}

	cfg := reformatConfig{
		filePath: c.Args().Get(0),
		date:     date,
		name:     c.String("name"),
		types:    c.String("types"),
	}

	err = validate(cfg)
	if err != nil {
		return err
	}

	err = Reformat(cfg)
	if err != nil {
		l.Fatalf("[INVOICE][REFACTOR] %+v", err)
		return err
	}

	return nil
}

func Reformat(c reformatConfig) error {
	name := generateName(c.name, c.types, c.date)

	err := rename(c.filePath, name)
	if err != nil {
		return errors.Wrap(err, "Can't reformat because of rename error")
	}

	return nil
}

func validate(c reformatConfig) error {
	if c.filePath == "" {
		return errors.New("FILE_PATH is required, given " + c.filePath)
	}

	_, err := os.Stat(c.filePath)
	if err != nil {
		return err
	}

	if c.name == "" {
		return errors.New("`NAME` is required")
	}
	return nil
}

func generateName(name, types string, date time.Time) string {
	return fmt.Sprintf("%s-%s-%s.pdf", strcase.ToKebab(date.Format("2006 01")), strcase.ToKebab(name), strcase.ToKebab(types))
}

func rename(filePath, newName string) error {
	dir := path.Dir(filePath)
	newPath := path.Join(dir, newName)
	err := os.Rename(filePath, newPath)
	l.Infof("File renamed to %s", newPath)
	return err
}
