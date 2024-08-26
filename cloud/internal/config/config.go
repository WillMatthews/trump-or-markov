package config

import (
	"fmt"
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
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s Server) Address() string {
	port := fmt.Sprintf("%d", s.Port)
	return s.Host + ":" + port
}

type TrumpTwitter struct {
	MaxTweets int    `yaml:"max_tweets"`
	Markov    Markov `yaml:"markov"`

	DoubleSpaceProb float64 `yaml:"double_space_prob"`
}

type Markov struct {
	MaxOrder            int `yaml:"max_order"`
	MaxGenerateAttempts int `yaml:"max_generate_attempts"`
	MaxChars            int `yaml:"max_chars"`
	MinWords            int `yaml:"min_words"`
	MaxWords            int `yaml:"max_words"` // TODO impl for general case

	EndPunctuation     []string `yaml:"end_punctuation"`
	EndPunctuationProb float64  `yaml:"end_punctuation_prob"`
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
	TrumpTwitter TrumpTwitter `yaml:"trump_twitter"`
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
