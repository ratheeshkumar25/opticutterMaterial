package utils

import (
	"errors"
	"fmt"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
)

// PredefinedSizes holds predefined sizes (fixed size ID => size)
var PredefinedSizes = map[uint]model.PredefinedSize{
	1: {Length: 50, Width: 30, Name: "Small"},
	2: {Length: 100, Width: 60, Name: "Medium"},
	3: {Length: 150, Width: 90, Name: "Large"},
}

// CalculateEstPrice calculates the estimated price based on item size.
func CalculateEstPrice(item *model.Item, pricePerUnit float64) (float64, error) {
	var length, width uint

	// Check if the item has a predefined size
	if item.FixedSizeID != 0 {
		// Fetch predefined size from the map
		size, exists := PredefinedSizes[item.FixedSizeID]
		if !exists {
			return 0, fmt.Errorf("invalid FixedSizeID: %d", item.FixedSizeID)
		}
		length = size.Length
		width = size.Width
	} else {
		// Use custom dimensions
		length = item.Length
		width = item.Width
	}

	if length == 0 || width == 0 {
		return 0, errors.New("invalid dimensions")
	}

	// Calculate estimated price (assuming it's based on area: length * width * price per unit)
	estPrice := float64(length*width) * pricePerUnit
	return estPrice, nil
}
