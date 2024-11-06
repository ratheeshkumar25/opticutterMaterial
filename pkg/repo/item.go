package repo

import "github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"

// CreateItem implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) CreateItem(Item *model.Item) (uint, error) {
	if err := m.DB.Create(&Item).Error; err != nil {
		return 0, err
	}
	return Item.ID, nil
}

// FindItemByID implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindItemByID(ItemsID uint) (*model.Item, error) {
	var Item model.Item

	if err := m.DB.First(&Item, ItemsID).Error; err != nil {
		return nil, err
	}
	return &Item, nil
}

// FindAllItem implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindAllItem() (*[]model.Item, error) {
	var Item []model.Item
	if err := m.DB.Find(&Item).Error; err != nil {
		return nil, err
	}
	return &Item, nil
}

// UpdateItem implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdateItem(Item *model.Item) error {
	if err := m.DB.Save(&Item).Error; err != nil {
		return err
	}
	return nil
}

// DeletItem implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) DeletItem(ItemID uint) error {
	if err := m.DB.Delete(&model.Item{}, ItemID).Error; err != nil {
		return err
	}
	return nil
}

// Get All items by user from database
func (m *MaterialRepository) FindAllItemByUsers(userID uint) (*[]model.Item, error) {
	var items []model.Item
	// Filter items by user ID
	if err := m.DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return &items, nil
}
