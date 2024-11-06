package service

import (
	"fmt"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

// PlaceOrderService implements interfaces.MaterialServiceInte.
func (m *MaterialService) PlaceOrderService(p *pb.Order) (*pb.OrderResponse, error) {
	item, err := m.Repo.FindItemByID(uint(p.Item_ID))
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "item not found",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, err
	}

	// Calculate totalPrice based on quantity
	totalAmount := item.EstPrice * float32(p.Quantity)

	newOrder := model.Order{
		UserID:    uint(p.User_ID),
		ItemID:    uint(p.Item_ID),
		Quantity:  int(p.Quantity),
		Status:    "Pending",
		CustomCut: p.CustomCut,
		IsCustom:  p.Is_Custom,
		Amount:    float64(totalAmount),
		PaymentID: p.Payment_ID,
	}

	// Save the order to the database
	orderID, err := m.Repo.CreateOrders(&newOrder)
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "failed to place order",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, err
	}

	// Update the material stock after placing the order
	err = m.Repo.UpdateMaterialStock(item.MaterialID, uint(p.Quantity))
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "failed to update material stock",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, nil
	}

	return &pb.OrderResponse{
		Status:  pb.OrderResponse_OK,
		Message: "Order placed successfully",
		Payload: &pb.OrderResponse_Data{
			Data: fmt.Sprintf("OrderID:%d, Total Amount: %.2f", orderID, totalAmount),
		},
	}, nil
}

// func (m *MaterialService) PlaceOrderService(p *pb.Order) (*pb.OrderResponse, error) {
// 	item, err := m.Repo.FindItemByID(uint(p.Item_ID))
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "item not found",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	//calculate totalPrice based on quantity
// 	totalAmount := item.EstPrice * float32(p.Quantity)

// 	newOrder := model.Order{
// 		UserID:    uint(p.User_ID),
// 		ItemID:    uint(p.Item_ID),
// 		Quantity:  int(p.Quantity),
// 		Status:    "Pending",
// 		CustomCut: p.CustomCut,
// 		IsCustom:  p.Is_Custom,
// 		Amount:    float64(totalAmount),
// 		PaymentID: p.Payment_ID,
// 	}

// 	//save the order to the database
// 	orderID, err := m.Repo.CreateOrders(&newOrder)
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "item not found",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	return &pb.OrderResponse{
// 		Status:  pb.OrderResponse_OK,
// 		Message: "Order placed successfully",
// 		Payload: &pb.OrderResponse_Data{
// 			Data: fmt.Sprintf("OrderID:%d, Total Amount: %.2f", orderID, totalAmount),
// 		},
// 	}, nil

// }

// FindAllOrdersSvc implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindAllOrdersSvc(p *pb.ItemNoParams) (*pb.OrderList, error) {
	result, err := m.Repo.FindAllOrders()
	if err != nil {
		return nil, err
	}

	var orders []*pb.Order
	for _, order := range *result {

		orders = append(orders, &pb.Order{
			Order_ID:   uint32(order.ID),
			User_ID:    uint32(order.UserID),
			Item_ID:    uint32(order.ItemID),
			Status:     order.Status,
			Amount:     order.Amount,
			Is_Custom:  order.IsCustom,
			Payment_ID: order.PaymentID,
		})
	}

	return &pb.OrderList{
		Orders: orders,
	}, nil
}

// FindOrderSvc implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindOrderSvc(p *pb.ItemID) (*pb.Order, error) {
	order, err := m.Repo.FindOrdersByID(uint(p.ID))
	if err != nil {
		return nil, err
	}

	return &pb.Order{
		Order_ID:   uint32(order.ID),
		User_ID:    uint32(order.UserID),
		Item_ID:    uint32(order.ItemID),
		Status:     order.Status,
		Amount:     order.Amount,
		Is_Custom:  order.IsCustom,
		Payment_ID: order.PaymentID,
	}, nil
}

// FindOrdersByUserSvc implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindOrdersByUserSvc(p *pb.ItemID) (*pb.OrderList, error) {
	userID := p.ID
	result, err := m.Repo.FindOrdersByUser(uint(userID))
	if err != nil {
		return nil, err
	}

	//var orders []*pb.Order
	ordList := &pb.OrderList{
		Orders: make([]*pb.Order, 0, len(*result)),
	}
	for _, order := range *result {

		ordList.Orders = append(ordList.Orders, &pb.Order{
			Order_ID:   uint32(order.ID),
			User_ID:    uint32(order.UserID),
			Item_ID:    uint32(order.ItemID),
			Status:     order.Status,
			Amount:     order.Amount,
			Payment_ID: order.PaymentID,
		})
	}

	return ordList, nil
}
