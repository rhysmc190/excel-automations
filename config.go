package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Directory string `yaml:"directory"`
	LogLevel  string `yaml:"logLevel"`
}

func loadConfig() Config {
	var cfg Config

	f, err := os.Open("config.yaml")
	processErrorWithoutLogger(err, "Error reading config file")
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	processErrorWithoutLogger(err, "Error parsing config file")

	return cfg
}
