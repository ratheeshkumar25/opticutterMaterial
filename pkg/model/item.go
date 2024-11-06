package model

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ItemName    string  `gorm:"not null"`
	MaterialID  uint    `gorm:"not null"`
	Length      uint    `gorm:"null"`
	Width       uint    `gorm:"null"`
	FixedSizeID uint    `gorm:"null"`
	IsCustom    bool    `gorm:"default:false"` // Flag indicating if size is custom
	EstPrice    float32 `gorm:"null"`
	UserID      uint    `gorm:"null"`
}

type PredefinedSize struct {
	gorm.Model
	Length uint   `gorm:"not null"`
	Width  uint   `gorm:"not null"`
	Name   string `gorm:"not null"` // Optional: Name like "Small", "Medium", "Large"
}
