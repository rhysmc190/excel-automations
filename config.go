package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputDirectory  string `yaml:"inputDirectory"`
	LogLevel        string `yaml:"logLevel"`
	OutputDirectory string `yaml:"outputDirectory"`
}

func getAbsPath(directory string) string {
	if directory == "" {
		return ""
	}
	absPath, err := filepath.Abs(directory)
	processError(err)
	return absPath
}

func loadConfig() Config {
	var cfg Config

	f, err := os.Open("config.yaml")
	processErrorWithoutLogger(err, "Error reading config file")
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	processErrorWithoutLogger(err, "Error parsing config file")

	cfg.InputDirectory = getAbsPath(cfg.InputDirectory)
	cfg.OutputDirectory = getAbsPath(cfg.OutputDirectory)

	return cfg
}
