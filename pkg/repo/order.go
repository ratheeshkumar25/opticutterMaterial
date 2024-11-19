package repo

import (
	"fmt"
	"log"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	"gorm.io/gorm"
)

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

// UpdateOrderStaus implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdateOrderStaus(OrderID uint, status string) error {
	// Find the order by ID and update the status
	if err := m.DB.Model(&model.Order{}).
		Where("id = ?", OrderID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}
	log.Println("Order status updated successfully")
	return nil
}

// GetLatestPaymentByOrderID implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) GetLatestPaymentByOrderID(orderID int) (model.Payment, error) {
	var payment model.Payment
	query := `
        SELECT payment_id, order_id, amount, status, client_secret, payment_method, user_id
        FROM payments
        WHERE order_id = ?
        ORDER BY 
            CASE 
                WHEN status = 'Pending' THEN 1
                WHEN status = 'Completed' THEN 2
                ELSE 3
            END, 
            payment_id DESC
        LIMIT 1`
	result := m.DB.Raw(query, orderID).Scan(&payment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return payment, nil // No payment found for the order_id
		}
		return payment, fmt.Errorf("failed to fetch latest payment: %v", result.Error)
	}
	return payment, nil
}
