package handlers

import (
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
	inter "github.com/ratheeshkumar25/opt_cut_material_service/pkg/service/interfaces"
)

// MaterialHandler represents the service methods and gRPC server for material-related operations.
type MaterialHandler struct {
	SVC inter.MaterialServiceInte
	pb.MaterialServiceServer
}

// NewMaterialHandler creates a new instance of productHandler with the provided product service interface.
func NewMaterialHandler(svc inter.MaterialServiceInte) *MaterialHandler {
	return &MaterialHandler{
		SVC: svc,
	}
}
