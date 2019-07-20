package gocreate

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

type config struct {
	Name          string
	directoryPath string
	force         string
}

type fileConfig struct {
	path         string
	templateName string
	templatePath string
}

type directoryConfig struct {
	path string
}

// Handler holds all logic about create subcommand
func Handler(c *cli.Context) error {
	cfg := config{
		Name:          "psutils",
		directoryPath: c.Args().Get(1),
		force:         c.Args().Get(0),
	}
	start(cfg)
	return nil
}

func start(cfg config) {
	log.Println("[OK] Starting...")

	cfg.directoryPath = normalizeDirectoryPath(cfg.directoryPath)
	var err error

	if cfg.force == "true" {
		err = cleanupDirectory(cfg.directoryPath)
		if err != nil {
			log.Fatalf("[ERROR] cleaning up directory %+v", err)
			return
		}
	}

	directories := []directoryConfig{
		{
			path: "",
		}, {
			path: "pkg",
		},
	}
	for _, directory := range directories {
		err = createDirectory(cfg.directoryPath + directory.path)
		if err != nil {
			log.Fatalf("[ERROR] creating directory %s: %+v", directory.path, err)
			return
		}
	}

	files := []fileConfig{{
		path:         "Makefile",
		templateName: "makefile",
		templatePath: "pkg/go-create/templates/makefile",
	}, {
		path:         "main.go",
		templateName: "main",
		templatePath: "pkg/go-create/templates/main",
	}, {
		path:         ".realize.yaml",
		templateName: "realize",
		templatePath: "pkg/go-create/templates/realize",
	}}

	for _, file := range files {
		err = createFile(file.templateName, filepath.Join(cfg.directoryPath, file.path), file.templatePath, cfg)
		if err != nil {
			log.Fatalf("[ERROR] creating %s: %+v", file.path, err)
			return
		}
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	exPath = normalizeDirectoryPath(exPath)
	log.Printf("[OK] Project %s created at %s\n", cfg.Name, filepath.Join(exPath, cfg.directoryPath))
}

func normalizeDirectoryPath(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}
	return path + "/"
}

func createDirectory(path string) error {
	_, err := ioutil.ReadDir(path)
	if err == nil {
		return errors.New("Directory is not empty")
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err == nil {
		return err
	}

	return nil
}

func cleanupDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil && err.Error() != "remove "+path+": no such file or directory" {
		return err
	}
	return nil
}

func createFile(name, path, templatePath string, cfg config) error {
	t := template.Must(template.New(name).ParseFiles(templatePath))
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = t.Execute(file, cfg)
	if err != nil {
		return err
	}
	return nil
}
