package repo

import "github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"

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
	if err := m.DB.Delete(model.Material{}, MaterialID).Error; err != nil {
		return err
	}
	return nil
}