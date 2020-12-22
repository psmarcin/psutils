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
	date                time.Time
	rootDirPath         string
	destinationFilePath string
	sourcePath          string
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

	rootDestinationDirPath := generateDirectoryPath(cfg.rootDirPath, cfg.date)
	_, fileName := path.Split(cfg.sourcePath)
	cfg.destinationFilePath = path.Join(rootDestinationDirPath, fileName)
	destinationDir, _ := path.Split(cfg.destinationFilePath)

	err = createDestinationDir(destinationDir)
	if err != nil {
		return errors.Wrap(err, "can't create destination directory")
	}

	err = cp(cfg.sourcePath, cfg.destinationFilePath)
	if err != nil {
		return errors.Wrap(err, "can't copy file")
	}

	logrus.Infof("File moved to %s", cfg.destinationFilePath)
	return nil
}

func generateDirectoryPath(root string, date time.Time) string {
	dir := path.Dir(root)
	dirWithDate := path.Join(dir, date.Format("2006 01"))

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

	_, err = os.Stat(path.Join(c.destinationFilePath, c.sourcePath))
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
