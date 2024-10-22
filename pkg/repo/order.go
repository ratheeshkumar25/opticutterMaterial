package repo

import "github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"

// CreateOrders implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) CreateOrders(order *model.Order) (uint, error) {
	if err := m.DB.Create(&order).Error; err != nil {
		return 0, err
	}
	return order.ID, nil
}

// FindAllOrders implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindAllOrders() (*[]model.Order, error) {
	var orders []model.Order
	if err := m.DB.Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

// FindOrder implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindOrder(userID uint, itemID uint) (*model.Order, error) {
	var order model.Order
	if err := m.DB.Where("user_id=? AND item_id=?", userID, itemID).Find(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// FindOrdersByID implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindOrdersByID(OrdersID uint) (*model.Order, error) {
	var Orders model.Order
	if err := m.DB.First(&Orders, OrdersID).Error; err != nil {
		return nil, err
	}
	return &Orders, nil
}

// FindOrdersByUser implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) FindOrdersByUser(userID uint) (*[]model.Order, error) {
	var orders []model.Order
	if err := m.DB.Where("user_id=?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

// UpdateOrders implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdateOrders(Orders *model.Order) error {
	if err := m.DB.Save(&Orders).Error; err != nil {
		return err
	}
	return nil
}

// DeleteOrders implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) DeleteOrders(OrdersID uint) error {
	if err := m.DB.Delete(&model.Order{}, OrdersID).Error; err != nil {
		return err
	}
	return nil
}
