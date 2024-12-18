package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func createCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
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

func writeEstimatesToCSVFile(estimates [][]string, filename string) {
	index := strings.Index(filename, ".xlsx")
	outfile := filename[:index] + "_BuildExactFormatted.csv"
	logger.Info().
		Str("outfile", outfile).
		Int("rows", len(estimates)).
		Msg("Writing data to outfile")
	writer, file, err := createCSVWriter(outfile)
	if err != nil {
		fmt.Println("Error creating CSV writer:", err)
		return
	}
	defer file.Close()
	for _, estimate := range estimates {
		writeCSVRecord(writer, estimate)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing CSV writer:", err)
	}
}
