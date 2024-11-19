package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	ItemID    uint    `gorm:"not null"`      // The item being ordered
	Quantity  int     `gorm:"not null"`      // Order quantity
	Status    string  `gorm:"not null"`      // Order status (Pending, Completed, etc.)
	CustomCut string  `gorm:"not null"`      // Custom cut or any other preferences
	IsCustom  bool    `gorm:"default:false"` // If the item has a custom size
	Amount    float64 `gorm:"not null"`
	//PaymentID string  `gorm:"not null"`
}
