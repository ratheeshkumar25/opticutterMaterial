package model

import "gorm.io/gorm"

// Component represents each panel piece that will be cut
type Component struct {
	gorm.Model
	MaterialID      uint   `gorm:"not null"`
	DoorPanel       string `gorm:"type:text"`
	BackSidePanel   string `gorm:"type:text"`
	SidePanel       string `gorm:"type:text"`
	TopPanel        string `gorm:"type:text"`
	BottomPanel     string `gorm:"type:text"`
	ShelvesPanel    string `gorm:"type:text"`
	PanelCount      int32  `gorm:"not null;default:0"`
	CuttingResultID uint   `gorm:"index"`
}

// CuttingResult represents the cutting result associated with an Item
type CuttingResult struct {
	gorm.Model
	ItemID     uint        `gorm:"not null;index"`
	Components []Component `gorm:"foreignKey:CuttingResultID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
