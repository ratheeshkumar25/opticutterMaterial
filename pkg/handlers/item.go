package handlers

import (
	"context"

	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

func (m *MaterialHandler) AddItem(ctx context.Context, p *pb.Item) (*pb.ItemResponse, error) {
	response, err := m.SVC.AddItemService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) FindItemByID(ctx context.Context, p *pb.ItemID) (*pb.Item, error) {
	response, err := m.SVC.FindItemByID(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) FindAllItem(ctx context.Context, p *pb.ItemNoParams) (*pb.ItemList, error) {
	response, err := m.SVC.FindAllItem(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) EditItem(ctx context.Context, p *pb.Item) (*pb.Item, error) {
	response, err := m.SVC.EditItemService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (m *MaterialHandler) RemoveItem(ctx context.Context, p *pb.ItemID) (*pb.ItemResponse, error) {
	response, err := m.SVC.RemoveItemService(p)
	if err != nil {
		return response, err
	}
	return response, nil
}
