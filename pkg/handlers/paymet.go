package handlers

import (
	"context"

	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

func (m *MaterialHandler) CreatePayment(ctx context.Context, p *pb.Order) (*pb.PaymentResponse, error) {
	response, err := m.SVC.PaymentService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (m *MaterialHandler) PaymentSuccess(ctx context.Context, p *pb.Payment) (*pb.PaymentStatusResponse, error) {
	response, err := m.SVC.PaymentSuccessService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
