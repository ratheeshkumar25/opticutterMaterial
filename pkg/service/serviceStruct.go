package service

import (
	inter "github.com/ratheeshkumar25/opt_cut_material_service/pkg/repo/interfaces"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/service/interfaces"
)

type MaterialService struct {
	Repo inter.MaterialRepoInter
}

func NewMaterialService(repo inter.MaterialRepoInter) interfaces.MaterialServiceInte {
	return &MaterialService{
		Repo: repo,
	}
}
