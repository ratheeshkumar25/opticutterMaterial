package utils

import (
	"errors"
	"fmt"
	"math"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
)

// PredefinedSizes holds predefined sizes (fixed size ID => size)
var PredefinedSizes = map[uint]model.PredefinedSize{
	1: {Length: 50, Width: 30, Name: "Small"},
	2: {Length: 100, Width: 60, Name: "Medium"},
	3: {Length: 150, Width: 90, Name: "Large"},
}

func CalculateEstPrice(item *model.Item, pricePerSheet, plywoodSize, wasteFactor float64) (float64, error) {
	// Calculate the required sheets using the helper function
	requiredSheets, err := CalculateRequiredSheets(item, plywoodSize, wasteFactor)
	if err != nil {
		return 0, err
	}

	// Calculate estimated price based on the number of sheets
	estPrice := requiredSheets * pricePerSheet
	roundedEstPrice := math.Round(estPrice*100) / 100

	return roundedEstPrice, nil
}

// CalculateRequiredSheets calculates the number of plywood sheets required based on item area, plywood sheet size, and waste factor.
func CalculateRequiredSheets(item *model.Item, plywoodSize, wasteFactor float64) (float64, error) {
	var length, width uint

	// Determine the dimensions of the item
	if item.FixedSizeID != 0 {
		size, exists := PredefinedSizes[item.FixedSizeID]
		if !exists {
			return 0, fmt.Errorf("invalid FixedSizeID: %d", item.FixedSizeID)
		}
		length = size.Length
		width = size.Width
	} else {
		length = item.Length
		width = item.Width
	}

	if length == 0 || width == 0 {
		return 0, errors.New("invalid dimensions")
	}

	// Calculate the item area
	itemArea := float64(length * width)

	// Calculate the required number of sheets, considering the waste factor
	requiredSheets := (itemArea * (1 + wasteFactor)) / plywoodSize
	return requiredSheets, nil
}
