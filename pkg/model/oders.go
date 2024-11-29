package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	ItemID    uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Status    string  `gorm:"not null"`
	CustomCut string  `gorm:"not null"`
	IsCustom  bool    `gorm:"default:false"`
	Amount    float64 `gorm:"not null"`
	Email     string  `gorm:"not null"`
	//PaymentID string  `gorm:"not null"`
}
