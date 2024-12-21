package utils

import (
	"fmt"
	"log"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
)

// Available Materials List (can be dynamically extended with more materials)
var materials = []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // Example list of material IDs

// Possible Item Names (Independent of material selection)
var itemNames = []string{
	"Shoe Rack",
	"Wardrobe",
}

// Component Generator for each Item Name
var ItemGenerators = map[string]func(materialID, length, width uint) []model.Component{
	"Shoe Rack": generateShoeRackComponents,
	"Wardrobe":  generateWardrobeComponents,
}

// GenerateComponents generates components for the selected item name
func GenerateComponents(materialID uint, itemName string, length, width uint) ([]model.Component, error) {
	// Ensure the materialID exists in the available materials list
	if !isMaterialValid(materialID) {
		return nil, fmt.Errorf("material ID %d is not valid", materialID)
	}

	// Ensure the item name exists in the list of possible items
	if !isItemValid(itemName) {
		return nil, fmt.Errorf("item %s is not valid", itemName)
	}

	// Fetch the corresponding component generator for the item name
	generator, exists := ItemGenerators[itemName]
	if !exists {
		return nil, fmt.Errorf("no generator found for item %s", itemName)
	}

	// Generate the components based on the item name and materialID
	return generator(materialID, length, width), nil
}

// Helper function to check if a material ID is valid
func isMaterialValid(materialID uint) bool {
	for _, material := range materials {
		if material == materialID {
			return true
		}
	}
	return false
}

// Helper function to check if an item name is valid
func isItemValid(itemName string) bool {
	for _, item := range itemNames {
		if item == itemName {
			return true
		}
	}
	return false
}

// Component Calculation Helper: This dynamically calculates panels based on available material.
func calculatePanelSize(length, width uint, panelType string) string {
	switch panelType {
	case "Door":
		return fmt.Sprintf("Door Panel (length: %d, width: %d)", length, width)
	case "Back":
		return fmt.Sprintf("Back Panel (length: %d, width: %d)", length/2, width/2)
	case "Side":
		return fmt.Sprintf("Side Panel (length: %d, width: %d)", length/3, width/3)
	case "Top":
		return fmt.Sprintf("Top Panel (length: %d, width: %d)", length/3, width/3)
	default:
		return "Unknown Panel"
	}
}

// Generate components for Shoe Rack
func generateShoeRackComponents(materialID, length, width uint) []model.Component {
	log.Println("Generating components for Shoe Rack with material ID", materialID)
	components := []model.Component{
		{
			MaterialID:    materialID,
			DoorPanel:     calculatePanelSize(length, width, "Door"),
			BackSidePanel: calculatePanelSize(length, width, "Back"),
			SidePanel:     calculatePanelSize(length, width, "Side"),
			TopPanel:      calculatePanelSize(length, width, "Top"),
			PanelCount:    4,
		},
	}
	return components
}

// Generate components for Wardrobe
func generateWardrobeComponents(materialID, length, width uint) []model.Component {
	log.Println("Generating components for Wardrobe with material ID", materialID)
	components := []model.Component{
		{
			MaterialID:    materialID,
			DoorPanel:     calculatePanelSize(length, width, "Door"),
			BackSidePanel: calculatePanelSize(length, width, "Back"),
			SidePanel:     calculatePanelSize(length, width, "Side"),
			TopPanel:      calculatePanelSize(length, width, "Top"),
			PanelCount:    4,
		},
	}
	return components
}

// package utils

// import (
// 	"errors"
// 	"fmt"
// 	"log"

// 	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
// )

// // ComponentGenerator is a function type that generates components
// type ComponentGenerator func(length, width uint) []model.Component

// // Registry of component generators
// var ComponentGenerators = map[uint]ComponentGenerator{
// 	1: generateShoeRackComponents,
// 	2: generateWardrobe1Components,
// 	3: generatWardrobeComponents,
// }

// // GenerateComponents dynamically generates components for a given materialID
// func GenerateComponents(materialID, length, width uint) ([]model.Component, error) {
// 	// lookup time for the generator function
// 	log.Println("matertialllID", materialID)
// 	generator, exists := ComponentGenerators[materialID]
// 	if !exists {
// 		return nil, errors.New("no component generator found for the given MaterialID")
// 	}

// 	// Call the appropriate generator function
// 	return generator(length, width), nil
// }

// // Helper functions to generate components basedWardrobe type
// func generateShoeRackComponents(length, width uint) []model.Component {
// 	// Perform calculations once and reuse results
// 	lengthDiv2 := length / 2
// 	widthDiv2 := width / 2
// 	lengthDiv3 := length / 3
// 	widthDiv3 := width / 3
// 	lengthDiv4 := length / 4
// 	widthDiv4 := width / 4
// 	lengthDiv5 := length / 5
// 	widthDiv5 := width / 5
// 	lengthDiv6 := length / 6
// 	widthDiv6 := width / 6

// 	return []model.Component{
// 		{
// 			MaterialID:    1,
// 			DoorPanel:     fmt.Sprintf("Shoe Rack Door Panel (Length: %d, Width: %d)", length, width),
// 			BackSidePanel: fmt.Sprintf("Shoe Rack Back Panel (Length: %d, Width: %d)", lengthDiv2, widthDiv2),
// 			SidePanel:     fmt.Sprintf("Shoe Rack Side Panel (Length: %d, Width: %d)", lengthDiv3, widthDiv3),
// 			TopPanel:      fmt.Sprintf("Shoe Rack Top Panel (Length: %d, Width: %d)", lengthDiv4, widthDiv4),
// 			BottomPanel:   fmt.Sprintf("Shoe Rack Bottom Panel (Length: %d, Width: %d)", lengthDiv5, widthDiv5),
// 			ShelvesPanel:  fmt.Sprintf("Shoe Rack Shelves Panel (Length: %d, Width: %d)", lengthDiv6, widthDiv6),
// 			PanelCount:    3,
// 		},
// 	}
// }

// func generateWardrobe1Components(length, width uint) []model.Component {
// 	// Perform calculations once and reuse results
// 	lengthDiv2 := length / 2
// 	widthDiv2 := width / 2
// 	lengthDiv3 := length / 3
// 	widthDiv3 := width / 3
// 	lengthDiv4 := length / 4
// 	widthDiv4 := width / 4
// 	lengthDiv5 := length / 5
// 	widthDiv5 := width / 5

// 	return []model.Component{
// 		{
// 			MaterialID:    2,
// 			DoorPanel:     fmt.Sprintf("Wardrobe1 Door Panel (Length: %d, Width: %d)", length, width),
// 			BackSidePanel: fmt.Sprintf("Wardrobe1 Back Panel (Length: %d, Width: %d)", length, widthDiv2),
// 			SidePanel:     fmt.Sprintf("Wardrobe1 Side Panel (Length: %d, Width: %d)", lengthDiv2, width),
// 			TopPanel:      fmt.Sprintf("Wardrobe1 Top Panel (Length: %d, Width: %d)", lengthDiv3, widthDiv3),
// 			BottomPanel:   fmt.Sprintf("Wardrobe1 Bottom Panel (Length: %d, Width: %d)", lengthDiv4, widthDiv4),
// 			ShelvesPanel:  fmt.Sprintf("Wardrobe1 Shelves Panel (Length: %d, Width: %d)", lengthDiv5, widthDiv5),
// 			PanelCount:    5,
// 		},
// 	}
// }

// // Add this function to your utils/component_generator.go file
// func generatWardrobeComponents(length, width uint) []model.Component {
// 	// Perform calculations once and reuse results
// 	lengthDiv2 := length / 2
// 	widthDiv2 := width / 2
// 	lengthDiv3 := length / 3
// 	widthDiv3 := width / 3
// 	lengthDiv4 := length / 4
// 	widthDiv4 := width / 4
// 	lengthDiv5 := length / 5
// 	widthDiv5 := width / 5

// 	return []model.Component{
// 		{
// 			MaterialID:    3,
// 			DoorPanel:     fmt.Sprintf("Wardrobe Door Panel (Length: %d, Width: %d)", length, width),
// 			BackSidePanel: fmt.Sprintf("Wardrobe Back Panel (Length: %d, Width: %d)", lengthDiv2, widthDiv2),
// 			SidePanel:     fmt.Sprintf("Wardrobe Side Panel (Length: %d, Width: %d)", lengthDiv3, widthDiv3),
// 			TopPanel:      fmt.Sprintf("Wardrobe Top Panel (Length: %d, Width: %d)", lengthDiv4, widthDiv4),
// 			BottomPanel:   fmt.Sprintf("Wardrobe Bottom Panel (Length: %d, Width: %d)", lengthDiv5, widthDiv5),
// 			ShelvesPanel:  fmt.Sprintf("Wardrobe Shelves Panel (Length: %d, Width: %d)", lengthDiv5, widthDiv5),
// 			PanelCount:    4,
// 		},
// 	}
// }
