package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

const QUANTITY = "Quantity"
const UNIT = "Unit"
const SUPPLY_RATE = "Supply Rate"
const LBR_RATE = "Lbr Rate"
const CATEGORY = "Category"
const DESCRIPTION = "Description"
const UOM = "UOM"
const UNIT_COST = "UnitCost"
const ITEM_TYPE = "ItemType"

const SHEET_NAME = "Standard Estimating"

var HEADING_NAMES = [4]string{QUANTITY, UNIT, SUPPLY_RATE, LBR_RATE}

func cleanString(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}

func getCleanedRowVals(r *xlsx.Row) ([]string, error) {
	var cleanedRowVals []string
	err := r.ForEachCell(func(c *xlsx.Cell) error {
		value, err := c.FormattedValue()
		cleanedRowVals = append(cleanedRowVals, cleanString(value))
		return err
	})
	return cleanedRowVals, err
}

func checkForCategoryHeading(r *xlsx.Row) (bool, map[string]int) {
	headingNameIndices := make(map[string]int)

	cleanedRowVals, err := getCleanedRowVals(r)
	processError(err)

	for _, headingName := range HEADING_NAMES {
		index := slices.Index(cleanedRowVals, cleanString(headingName))
		if index == -1 {
			return false, nil
		}
		headingNameIndices[headingName] = index
	}

	return true, headingNameIndices
}

func getCategoryFromHeadingRow(r *xlsx.Row) string {
	var category string
	err := r.ForEachCell(func(c *xlsx.Cell) error {
		val, err := c.FormattedValue()
		if category == "" && val != "" {
			category = strings.TrimSpace(val)
		}
		return err
	})
	processError(err)
	return category
}

func getDescriptionFromRowVals(rowVals []string) string {
	for _, val := range rowVals {
		if val != "" {
			return strings.TrimSpace(val)
		}
	}
	return ""
}

func getRoundedRate(rate string) float64 {
	rateAsFloat, err := strconv.ParseFloat(strings.ReplaceAll(rate, "$", ""), 64)
	if err != nil && rate != "" {
		fmt.Println(err)
		panic(errors.New("cannot parse rate"))
	}
	// round to 2 dp
	return math.Round(rateAsFloat*100) / 100
}

func getRoundedQuantity(quantity string) string {
	quantAsFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil && quantity != "" {
		fmt.Println(err)
		panic(errors.New("cannot parse quantity"))
	}
	return fmt.Sprintf("%.1f", quantAsFloat)
}

func parseSupplyAndLabourRates(rowVals []string, headingNameIndices map[string]int) (string, string) {
	var itemType string
	supplyRate := getRoundedRate(rowVals[headingNameIndices[SUPPLY_RATE]])
	lbrRate := getRoundedRate(rowVals[headingNameIndices[LBR_RATE]])

	if supplyRate != 0 {
		itemType = "Material"
	}
	if lbrRate != 0 {
		itemType = "Labour"
	}

	return strconv.FormatFloat(supplyRate+lbrRate, 'f', 2, 64), itemType
}

func parseDataRow(r *xlsx.Row, headingNameIndices map[string]int) (bool, string, string, string, string, string) {
	var rowVals []string
	err := r.ForEachCell(func(c *xlsx.Cell) error {
		value, err := c.FormattedValue()
		rowVals = append(rowVals, value)
		return err
	})
	processError(err)

	logger.Debug().Strs("rowVals", rowVals).Msg("Read row")

	description := getDescriptionFromRowVals(rowVals)
	quantity := getRoundedQuantity(rowVals[headingNameIndices[QUANTITY]])
	uom := strings.TrimSpace(rowVals[headingNameIndices[UNIT]])
	unitCost, itemType := parseSupplyAndLabourRates(rowVals, headingNameIndices)

	isDataRow := description != "" && quantity != "" && uom != "" && unitCost != ""

	logger.Debug().
		Bool("isDataRow", isDataRow).
		Str("description", description).
		Str("quantity", quantity).
		Str("uom", uom).
		Str("unitCost", unitCost).
		Str("itemType", itemType).
		Msg("Parsed row data")

	return isDataRow, description, quantity, uom, unitCost, itemType
}

func parseRows(s *xlsx.Sheet) [][]string {
	estimates := [][]string{{CATEGORY, DESCRIPTION, QUANTITY, UOM, UNIT_COST, ITEM_TYPE}}
	var category string
	var headingNameIndices map[string]int
	logger.Info().Msg("Parsing file")
	err := s.ForEachRow(func(r *xlsx.Row) error {
		isCategoryHeading, indices := checkForCategoryHeading(r)

		if isCategoryHeading {
			category = getCategoryFromHeadingRow(r)
			headingNameIndices = indices
			logger.Debug().Str("category", category).Msg("Found Heading Row")
			return nil
		}

		isDataRow, description, quantity, uom, unitCost, itemType := parseDataRow(r, headingNameIndices)
		if category != "" && isDataRow {
			estimates = append(estimates, []string{category, description, quantity, uom, unitCost, itemType})
		}
		return nil
	})
	processError(err)
	return estimates
}

func readExcelFile(filename string) *xlsx.Sheet {
	wb, err := xlsx.OpenFile(filename)
	processError(err)
	s, ok := wb.Sheet["Standard Estimating"]
	if !ok {
		err := errors.New("sheet not found")
		logger.Error().Err(err)
		panic(err)
	}
	logger.Info().Msg("Successfully read excel file")
	return s
}
