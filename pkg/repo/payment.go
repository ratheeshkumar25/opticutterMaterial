package repo

import (
	"fmt"
	"log"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
)

// SavePayment implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) SavePayment(payment *model.Payment) error {
	if err := m.DB.Create(&payment).Error; err != nil {
		return fmt.Errorf("failed to save payment: %v", err)
	}
	log.Println("Payment saved successfully")
	return nil
}

// UpdatePaymentAndOrderStatus implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdatePaymentAndOrderStatus(paymentID string, orderID int, paymentStatus string, orderStatus string) error {
	// Begin a transaction
	tx := m.DB.Begin()

	// Update payment status
	if err := tx.Model(&model.Payment{}).
		Where("payment_id = ?", paymentID).
		Update("status", paymentStatus).Error; err != nil {
		tx.Rollback() // Roll back the transaction on error
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	// Update order status
	if err := tx.Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", orderStatus).Error; err != nil {
		tx.Rollback() // Roll back the transaction on error
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// Commit the transaction if both updates succeed
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("Payment and order status updated successfully")
	return nil
}

// UpdatePaymentStatus implements interfaces.MaterialRepoInter.
func (m *MaterialRepository) UpdatePaymentStatus(paymentID string, status string) error {
	// Find the payment by ID and update the status
	if err := m.DB.Model(&model.Payment{}).
		Where("id = ?", paymentID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}
	log.Println("Payment status updated successfully")
	return nil
}
