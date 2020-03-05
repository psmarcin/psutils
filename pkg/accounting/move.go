package accounting

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"psutils/pkg/config"
)

type moveConfig struct {
	date               time.Time
	rootDirPath        string
	destinationDirPath string
	destinationPath    string
	sourcePath         string
}

// Handler holds all logic about create subcommand
func MoveHandler(c *cli.Context) error {
	conf := config.Load()
	date, err := time.Parse(conf.Other.MontDateFormat, c.String("date"))
	if err != nil {
		return err
	}

	cfg := moveConfig{
		sourcePath:  c.Args().Get(0),
		date:        date,
		rootDirPath: conf.Accounting.FilesDirectory,
	}

	err = moveValidate(cfg)
	if err != nil {
		return errors.Wrap(err, "validate source and destination")
	}

	dir := generateDirectoryPath(cfg)
	cfg.destinationPath = path.Join(dir, cfg.sourcePath)
	cfg.destinationDirPath, _ = path.Split(cfg.destinationPath)

	err = createDestinationDir(cfg.destinationDirPath)
	if err != nil {
		return errors.Wrap(err, "can't create destination directory")
	}

	err = cp(cfg.sourcePath, path.Join(dir, cfg.sourcePath))
	if err != nil {
		return errors.Wrap(err, "can't copy file")
	}

	logrus.Infof("File moved to %s", cfg.destinationPath)
	return nil
}

func generateDirectoryPath(cfg moveConfig) string {
	mainDir := path.Dir(cfg.destinationDirPath)
	dirWithDate := path.Join(mainDir, cfg.date.Format("2006 01"))

	return dirWithDate
}

func createDestinationDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func moveValidate(c moveConfig) error {
	if c.sourcePath == "" {
		return errors.New("FILE_PATH is required, given " + c.sourcePath)
	}

	_, err := os.Stat(c.sourcePath)
	if err != nil {
		return err
	}

	_, err = os.Stat(path.Join(c.destinationDirPath, c.sourcePath))
	if err == nil {
		return err
	}

	return nil
}

func cp(source, direction string) error {
	sourceStream, _ := os.Open(source)
	writeStream, _ := os.Create(direction)
	_, err := io.Copy(writeStream, sourceStream)
	return err
}
