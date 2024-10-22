package repo

import (
	inter "github.com/ratheeshkumar25/opt_cut_material_service/pkg/repo/interfaces"
	"gorm.io/gorm"
)

// MaterialRepository for connecting db to ProductRepoInter methods
type MaterialRepository struct {
	DB *gorm.DB
}

// NewMaterialRepository creates an instance of Material repo
func NewMaterialRepository(db *gorm.DB) inter.MaterialRepoInter {
	return &MaterialRepository{
		DB: db,
	}
}
