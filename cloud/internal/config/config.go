package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	ConfigPath = "./config.yaml"
)

type App struct {
	Name string
}

type Server struct {
	host string
	port int
}

type Dataset struct {
	Trump    string
	MobyDick string
}

type Database struct {
	Sqlite string
}

type Config struct {
	App      App
	Server   Server
	Dataset  Dataset
	Database Database
}

func GetConfig() (*Config, string) {
	yamlFile, err := os.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}

	if len(yamlFile) == 0 {
		panic("config file is empty")
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// todo make this a mustgetenv once I've got docker working
	version := os.Getenv("TOM_VERSION")

	return &config, version
}
