package handlers

import (
	"context"

	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

// CreateMaterial handles the creation of a new product.
func (m *MaterialHandler) AddMaterial(ctx context.Context, p *pb.Material) (*pb.MaterialResponse, error) {
	response, err := m.SVC.AddMaterialService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

// FindMaterialByID retrieves a material by its ID.
func (m *MaterialHandler) FindMaterialByID(ctx context.Context, p *pb.MaterialID) (*pb.Material, error) {
	response, err := m.SVC.FindMaterialByIDService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

// FindAllMaterial retrieves all material.
func (m *MaterialHandler) FindAllMaterial(ctx context.Context, p *pb.MaterialNoParams) (*pb.MaterialList, error) {
	response, err := m.SVC.FindAllMaterialService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

// EditMaterial handles the editing of an existing material.
func (m *MaterialHandler) EditMaterial(ctx context.Context, p *pb.Material) (*pb.Material, error) {
	response, err := m.SVC.EditMaterialService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

// RemoveMaterial  handles the removal of a material.
func (m *MaterialHandler) RemoveMaterial(ctx context.Context, p *pb.MaterialID) (*pb.MaterialResponse, error) {
	response, err := m.SVC.RemoveMaterialService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
