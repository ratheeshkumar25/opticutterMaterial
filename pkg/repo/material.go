package repo

import (
	"fmt"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
)

// AddMaterial implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) AddMaterial(Material *model.Material) (uint, error) {
	if err := m.DB.Create(&Material).Error; err != nil {
		return 0, err
	}
	return Material.ID, nil
}

// FindMaterialByID implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindMaterialByID(MaterialID uint) (*model.Material, error) {
	var Material model.Material
	if err := m.DB.First(&Material, MaterialID).Error; err != nil {
		return nil, err
	}
	return &Material, nil
}

// FindAllMaterial implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindAllMaterial() (*[]model.Material, error) {
	var Material []model.Material
	if err := m.DB.Find(&Material).Error; err != nil {
		return nil, err
	}
	return &Material, nil
}

// UpdateMaterial implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdateMaterial(Material *model.Material) error {
	if err := m.DB.Save(&Material).Error; err != nil {
		return err
	}
	return nil
}

// DeleteMaterial implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) DeleteMaterial(MaterialID uint) error {
	if err := m.DB.Delete(&model.Material{}, MaterialID).Error; err != nil {
		return err
	}
	return nil
}

// UpdateMaterialStock implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdateMaterialStock(materialID uint, quantity uint) error {
	// Fetch the material from the database
	var material model.Material
	err := m.DB.First(&material, materialID).Error
	if err != nil {
		return err
	}

	// Check if there is enough stock available
	if material.Stock < int(quantity) {
		return fmt.Errorf("not enough stock available")
	}

	// Deduct the quantity from the stock
	material.Stock -= int(quantity) // Convert quantity to int

	// Save the updated material
	return m.DB.Save(&material).Error
}
