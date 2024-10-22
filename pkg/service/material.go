package service

import (
	"fmt"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

// AddMaterialService implements interfaces.MaterialServiceInte.
func (m *MaterialService) AddMaterialService(p *pb.Material) (*pb.MaterialResponse, error) {
	// Create a new material instance from the input
	newMaterial := model.Material{
		Name:        p.Material_Name,
		Description: p.Description,
		Stock:       int(p.Stock),
		Price:       p.Price,
	}

	// Insert the new material into the database through repository
	materialID, err := m.Repo.AddMaterial(&newMaterial)
	if err != nil {
		// Return an error response if insertion fails
		return &pb.MaterialResponse{
			Status:  pb.MaterialResponse_ERROR,
			Message: "Failed to add material",
			Payload: &pb.MaterialResponse_Error{
				Error: err.Error(),
			},
		}, err
	}

	// If insertion is successful, return a success response
	return &pb.MaterialResponse{
		Status:  pb.MaterialResponse_OK,
		Message: "Material added successfully",
		Payload: &pb.MaterialResponse_Data{
			Data: fmt.Sprintf("MaterialID: %d", materialID),
		},
	}, nil
}

// EditMaterialService implements interfaces.MaterialServiceInte.
func (m *MaterialService) EditMaterialService(p *pb.Material) (*pb.Material, error) {
	material, err := m.Repo.FindMaterialByID(uint(p.Material_ID))
	if err != nil {
		return nil, err
	}

	material.Name = p.Material_Name
	material.Description = p.Description
	material.Stock = int(p.Stock)
	material.Price = p.Price

	err = m.Repo.UpdateMaterial(material)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// FindMaterialByIDService implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindMaterialByIDService(p *pb.MaterialID) (*pb.Material, error) {
	material, err := m.Repo.FindMaterialByID(uint(p.ID))
	if err != nil {
		return nil, err
	}

	return &pb.Material{
		Material_Name: material.Name,
		Description:   material.Description,
		Stock:         int32(material.Stock),
		Price:         material.Price,
	}, nil
}

// FindAllMaterialService implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindAllMaterialService(p *pb.MaterialNoParams) (*pb.MaterialList, error) {
	// Fetch all materials from the repository
	result, err := m.Repo.FindAllMaterial()
	if err != nil {
		return nil, err
	}

	// Prepare a slice to hold the  materials
	var materials []*pb.Material

	// Loop through the result and map each material to the protobuf format
	for _, material := range *result {
		pbMaterial := &pb.Material{
			Material_ID:   uint32(material.ID),
			Material_Name: material.Name,
			Description:   material.Description,
			Stock:         int32(material.Stock),
			Price:         material.Price,
		}
		materials = append(materials, pbMaterial)
	}

	return &pb.MaterialList{
		Materials: materials,
	}, nil
}

// RemoveMaterialService implements interfaces.MaterialServiceInte.
func (m *MaterialService) RemoveMaterialService(p *pb.MaterialID) (*pb.MaterialResponse, error) {
	// Find the material by its ID
	material, err := m.Repo.FindMaterialByID(uint(p.ID))
	if err != nil {
		return &pb.MaterialResponse{
			Status:  pb.MaterialResponse_ERROR,
			Message: "Error in finding material",
			Payload: &pb.MaterialResponse_Error{Error: err.Error()},
		}, err
	}

	// Remove the material from the repository
	err = m.Repo.DeleteMaterial(material.ID)
	if err != nil {
		return &pb.MaterialResponse{
			Status:  pb.MaterialResponse_ERROR,
			Message: "Error in removing material",
			Payload: &pb.MaterialResponse_Error{Error: err.Error()},
		}, err
	}

	// Return success response after successful removal
	return &pb.MaterialResponse{
		Status:  pb.MaterialResponse_OK,
		Message: "Material removed successfully",
	}, nil

}
