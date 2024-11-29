package handlers

import (
	"context"

	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

func (m *MaterialHandler) GetCuttingResult(ctx context.Context, p *pb.ItemID) (*pb.CuttingResultResponse, error) {
	response, err := m.SVC.GetCuttingResService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
