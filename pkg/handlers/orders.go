package handlers

import (
	"context"

	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

func (m *MaterialHandler) PlaceOrder(ctx context.Context, p *pb.Order) (*pb.OrderResponse, error) {
	response, err := m.SVC.PlaceOrderService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) OrderHistory(ctx context.Context, p *pb.ItemNoParams) (*pb.OrderList, error) {
	response, err := m.SVC.FindAllOrdersSvc(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) FindOrder(ctx context.Context, p *pb.ItemID) (*pb.Order, error) {
	response, err := m.SVC.FindOrderSvc(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) FindOrdersByUser(ctx context.Context, p *pb.ItemID) (*pb.OrderList, error) {
	response, err := m.SVC.FindOrdersByUserSvc(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
