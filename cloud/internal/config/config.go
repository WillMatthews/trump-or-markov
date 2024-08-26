package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	ConfigPath = "./config.yaml"
)

type App struct {
	Name string `yaml:"name"`
}

type Server struct {
	host string `yaml:"host"`
	port int    `yaml:"port"`
}

type Twitter struct {
	MaxTweets int    `yaml:"max_tweets"`
	Markov    Markov `yaml:"markov"`
}

type Markov struct {
	MaxOrder            int `yaml:"max_order"`
	MaxGenerateAttempts int `yaml:"max_generate_attempts"`
	MaxChars            int `yaml:"max_chars"`
	MinWords            int `yaml:"min_words"`
	MaxWords            int `yaml:"max_words"` // TODO impl for general case
}

type Dataset struct {
	Trump    string `yaml:"trump"`
	MobyDick string `yaml:"mobydick"`
}

type Database struct {
	Sqlite string `yaml:"sqlite"`
}

type Config struct {
	App     App     `yaml:"app"`
	Server  Server  `yaml:"server"`
	Dataset Dataset `yaml:"dataset"`
	//Database Database `yaml:"database"`
	Twitter Twitter `yaml:"twitter"`
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
