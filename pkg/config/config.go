package config

import (
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var l log.Logger
var configFilePath string
var configPath string

type Config struct {
	Version    string     `yaml:"version"`
	Accounting Accounting `yaml:"accounting"`
	Other      Other      `yaml:"other"`
}

type Accounting struct {
	FilesDirectory string `yaml:"filesDirectory"`
	Confirmation
}

type Company struct {
	Name     string `yaml:"name"`
	Address1 string `yaml:"address1"`
	Address2 string `yaml:"address2"`
	Nip      string `yaml:"nip"`
}

type Confirmation struct {
	Seller   Company  `yaml:"seller"`
	Customer Company  `yaml:"customer"`
	Items    []string `yaml:"items"`
}

type Other struct {
	MontDateFormat string `yaml:"month-date-format"`
}

func init() {
	l = log.With("package", "invoice")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.With("err", err).Fatalf("Can't get home directory")
	}

	configPath = path.Join(homeDir, ".psutils")

	configFilePath = path.Join(configPath, "config.yaml")
}

func Load() Config {
	config, err := openConfig()
	if err != nil {
		log.With("err", err).Infof("Can't open config file, creating default")
		err = createDefault()

		if err != nil {
			log.With("err", err).Fatalf("Can't create default config file")
		}
		return Load()
	}
	return config
}

func openConfig() (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func createDefault() error {
	err := createDefaultDir()
	if err != nil {
		return err
	}

	err = createDefaultFile()
	if err != nil {
		return err
	}
	return nil
}

func createDefaultDir() error {
	_, err := ioutil.ReadDir(configPath)
	if err == nil {
		return nil
	}
	log.With("err", err).With("configPath", configPath).Info("Directory not exists")
	err = os.Mkdir(configPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultFile() error {
	config := Config{
		Version: "v1",
		Accounting: Accounting{
			FilesDirectory: "~/Downloads",
			Confirmation: Confirmation{
				Seller: Company{
					Name:     "Company",
					Address1: "Deepest Hell 4",
					Address2: "Heaven 123-50000",
					Nip:      "111111111",
				},
				Customer: Company{
					Name:     "Company",
					Address1: "Forest Guy",
					Address2: "Fire Stone 178",
					Nip:      "9999999999"},
				Items: []string{"Documentation"},
			},
		},
		Other: Other{
			MontDateFormat: "2006/01",
		},
	}

	str, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFilePath, str, os.ModePerm) // todo: fix file permissions
	if err != nil {
		return err
	}
	return nil
}
