package repo

import (
	"log"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	"gorm.io/gorm"
)

// SaveCuttingResult implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) SaveCuttingResult(itemID uint, components []model.Component) error {
	// Create the CuttingResult model
	log.Println("savecuttingresult", itemID)
	cuttingResult := model.CuttingResult{
		ItemID:     itemID,
		Components: components,
	}

	// Save the CuttingResult to the database
	if err := m.DB.Create(&cuttingResult).Error; err != nil {
		return err
	}
	return nil
}

func (m *MaterialRepository) GetCuttingResultByItemID(itemID uint) ([]model.Component, error) {
	var components []model.Component

	// Use GORM to find components associated with the cutting result
	result := m.DB.Where("cutting_result_id IN (SELECT id FROM cutting_results WHERE item_id = ?)", itemID).Find(&components)

	// Check for errors
	if result.Error != nil {
		return nil, result.Error
	}

	// Check if no components were found
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return components, nil
}

// func (m *MaterialRepository) GetCuttingResultByItemID(itemID uint) ([]model.Component, error) {
// 	var components []model.Component

// 	// SQL query - adjust based on your database schema
// 	query := `SELECT material_id, door_panel, back_side_panel, side_panel, top_panel, bottom_panel, shelves_panel, panel_count
// 			  FROM components
// 			  WHERE item_id = ?`

// 	// Execute the query and scan the results into the components slice
// 	err := m.DB.Select(&components, query, itemID)
// 	if err != nil {
// 		return nil, err.Error
// 	}

// 	return components, nil
// }
