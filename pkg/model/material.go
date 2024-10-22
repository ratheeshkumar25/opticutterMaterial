package model

import "gorm.io/gorm"

type Material struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	Price       float64 `gorm:"not null"`
}
