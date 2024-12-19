package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(filename)
			logger.Info().
				Str("dir", dir).
				Msg("Couldn't write csv file into non-existent directory, creating directory and trying again")
			e := os.MkdirAll(dir, 0755)
			if e == nil {
				return createCSVWriter(filename)
			}
		}
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	return writer, f, nil
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		fmt.Println("Error writing record to CSV:", err)
	}
}

func removeFileExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func getOutfileName(filename string) string {
	var outfile string
	if config.OutputDirectory != "" {
		outputDir, err := filepath.Abs(config.OutputDirectory)
		processError(err)
		outfile = filepath.Join(outputDir, removeFileExtension(filepath.Base(filename)))
	} else {
		outfile = removeFileExtension(filename)
	}
	return outfile + "_BuildExactFormatted.csv"
}

func writeEstimatesToCSVFile(estimates [][]string, filename string) {
	outfile := getOutfileName(filename)
	writer, file, err := createCSVWriter(outfile)
	processError(err)
	defer file.Close()

	logger.Info().
		Str("outfile", outfile).
		Int("rows", len(estimates)).
		Msg("Writing data to outfile")
	for _, estimate := range estimates {
		writeCSVRecord(writer, estimate)
	}

	writer.Flush()
	err = writer.Error()
	processError(err)
}
