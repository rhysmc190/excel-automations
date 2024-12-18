package main

import (
	"os"
	"path/filepath"

	"github.com/akatsuki-py/tfd"
	"github.com/rs/zerolog"
)

func promptForFileName() string {
	directory, err := filepath.Abs(config.Directory)
	processError(err)
	os.Chdir(directory)
	filename, err := tfd.CreateSelectDialog([]string{"xlsx"}, false)
	processError(err)

	logger.Info().Str("filename", filename).Msg("Read in filename")

	return filename
}

var (
	config Config
	logger zerolog.Logger
)

func init() {
	config = loadConfig()
	logger = getLogger()
}

func main() {
	filename := promptForFileName()

	s := readExcelFile(filename)

	estimates := parseRows(s)

	writeEstimatesToCSVFile(estimates, filename)
}
