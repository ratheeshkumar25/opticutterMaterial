package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
)

func (m *MaterialService) PlaceOrderService(p *pb.Order) (*pb.OrderResponse, error) {

	item, err := m.Repo.FindItemByID(uint(p.Item_ID))
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "Item not found",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, err
	}

	// No need to check for existing order by Order_ID as it's a new order
	// Create a new order in the database
	totalAmount := item.EstPrice * float32(p.Quantity)

	newOrder := model.Order{
		UserID:    uint(p.User_ID),
		ItemID:    uint(p.Item_ID),
		Quantity:  int(p.Quantity),
		Status:    "Pending",
		CustomCut: p.CustomCut,
		IsCustom:  p.Is_Custom,
		Amount:    float64(totalAmount),
	}

	// Save the order to the database
	orderID, err := m.Repo.CreateOrders(&newOrder)
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "Failed to place order",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, err
	}

	// Update the material stock after placing the order
	err = m.Repo.UpdateMaterialStock(item.MaterialID, uint(p.Quantity))
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "Failed to update material stock",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, nil
	}

	// Prepare order details for Redis storage
	orderData := map[string]interface{}{
		"OrderID":       orderID, // Use the generated orderID
		"UserID":        newOrder.UserID,
		"PaymentAmount": newOrder.Amount,
		"Status":        newOrder.Status,
	}

	// Serialize order data to JSON
	orderDataJSON, err := json.Marshal(orderData)
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "Failed to serialize order data",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, err
	}

	// Define Redis key for order and store serialized data
	key := fmt.Sprintf("order:%d:user:%d", orderID, newOrder.UserID)

	// Store order details in Redis
	err = m.redis.SetDataInRedis(key, orderDataJSON, time.Hour) // setting expiration to 1 hour
	if err != nil {
		return &pb.OrderResponse{
			Status:  pb.OrderResponse_ERROR,
			Message: "Failed to store order data in Redis",
			Payload: &pb.OrderResponse_Error{Error: err.Error()},
		}, err
	}

	// Return success response with the new order details
	return &pb.OrderResponse{
		Status:  pb.OrderResponse_OK,
		Message: "Order placed successfully",
		Payload: &pb.OrderResponse_Data{
			Data: fmt.Sprintf("order_id:%d, total_amount: %.2f", orderID, totalAmount),
		},
	}, nil
}

// // PlaceOrderService implements interfaces.MaterialServiceInte.
// func (m *MaterialService) PlaceOrderService(p *pb.Order) (*pb.OrderResponse, error) {

// 	item, err := m.Repo.FindItemByID(uint(p.Item_ID))
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "item not found",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	// Check if the order has already been paid for
// 	order, err := m.Repo.FindOrdersByID(uint(p.Order_ID))
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "Order not found",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	// If the order status is 'Completed', notify the user that the payment was successful
// 	if order.Status == "Completed" {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "Payment has already been successfully processed for this order. Please contact customer service for further assistance.",
// 		}, nil
// 	}

// 	// Calculate totalPrice based on quantity
// 	totalAmount := item.EstPrice * float32(p.Quantity)

// 	// Generate a unique PaymentID (if it isn't provided by the caller)
// 	// paymentID := p.Payment_ID
// 	// if paymentID == "" {
// 	// 	// Generate a new UUID if the PaymentID is empty
// 	// 	paymentID = uuid.New().String()
// 	// }

// 	newOrder := model.Order{
// 		UserID:    uint(p.User_ID),
// 		ItemID:    uint(p.Item_ID),
// 		Quantity:  int(p.Quantity),
// 		Status:    "Pending",
// 		CustomCut: p.CustomCut,
// 		IsCustom:  p.Is_Custom,
// 		Amount:    float64(totalAmount),
// 		//PaymentID: paymentID,
// 	}

// 	// Save the order to the database
// 	orderID, err := m.Repo.CreateOrders(&newOrder)
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "failed to place order",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	// Update the material stock after placing the order
// 	err = m.Repo.UpdateMaterialStock(item.MaterialID, uint(p.Quantity))
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "failed to update material stock",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, nil
// 	}

// 	// Prepare order details for Redis storage
// 	orderData := map[string]interface{}{
// 		"OrderID":       newOrder.ID,
// 		"UserID":        newOrder.UserID,
// 		"PaymentAmount": newOrder.Amount,
// 		//"PaymentID":     newOrder.PaymentID,
// 		"Status": newOrder.Status,
// 	}

// 	// Serialize order data to JSON
// 	orderDataJSON, err := json.Marshal(orderData)
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "Failed to serialize order data",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	// Define Redis key for order and store serialized data
// 	key := fmt.Sprintf("order:%d:user:%d", orderID, newOrder.UserID)

// 	// redisKey := fmt.Sprintf("order:%d:user:%d", p.Order_ID, p.User_ID)
// 	fmt.Println("Storing order data in Redis with key:", key)

// 	// Store order details in Redis
// 	err = m.redis.SetDataInRedis(key, orderDataJSON, time.Hour) // setting expiration to 1 hour
// 	if err != nil {
// 		return &pb.OrderResponse{
// 			Status:  pb.OrderResponse_ERROR,
// 			Message: "Failed to store order data in Redis",
// 			Payload: &pb.OrderResponse_Error{Error: err.Error()},
// 		}, err
// 	}

// 	return &pb.OrderResponse{
// 		Status:  pb.OrderResponse_OK,
// 		Message: "Order placed successfully",
// 		Payload: &pb.OrderResponse_Data{
// 			Data: fmt.Sprintf("order_id:%d, total_amount: %.2f", orderID, totalAmount),
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
			Order_ID:  uint32(order.ID),
			User_ID:   uint32(order.UserID),
			Item_ID:   uint32(order.ItemID),
			Status:    order.Status,
			Amount:    order.Amount,
			Is_Custom: order.IsCustom,
			//Payment_ID: order.PaymentID,
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
		Order_ID:  uint32(order.ID),
		User_ID:   uint32(order.UserID),
		Item_ID:   uint32(order.ItemID),
		Status:    order.Status,
		Amount:    order.Amount,
		Is_Custom: order.IsCustom,
		//Payment_ID: order.PaymentID,
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
			Order_ID: uint32(order.ID),
			User_ID:  uint32(order.UserID),
			Item_ID:  uint32(order.ItemID),
			Status:   order.Status,
			Amount:   order.Amount,
			//Payment_ID: order.PaymentID,
		})
	}

	return ordList, nil
}
