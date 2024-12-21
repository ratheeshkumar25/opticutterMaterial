package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/utils"
)

// PaymentService handles payment intent creation
func (m *MaterialService) PaymentService(p *pb.Order) (*pb.PaymentResponse, error) {
	key := fmt.Sprintf("order:%d:user:%d", p.Order_ID, p.User_ID)
	log.Printf("Attempting to retrieve order data with key: %s", key)

	// Retrieve order data from Redis
	orderData, err := m.redis.GetFromRedis(key)
	if err != nil {
		if err == redis.Nil {
			log.Printf("No data found in Redis for key: %s", key)
			return nil, fmt.Errorf("no order data found for order ID %v and user ID %v", p.Order_ID, p.User_ID)
		}
		return nil, err
	}

	log.Printf("Order data retrieved from Redis for key %s: %s", key, orderData)

	var payment model.Payment
	err = json.Unmarshal([]byte(orderData), &payment)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order data: %v", err)
	}

	// Check if a PaymentIntent already exists and if it's incomplete
	if payment.PaymentID != "" && payment.Status != "Completed" {
		log.Printf("Existing PaymentIntent found with ID: %s", payment.PaymentID)
	} else {
		// If no PaymentIntent exists or if the status is "Completed," create a new one
		amountInCents := int(payment.Amount)
		if amountInCents < 100 { // Convert to cents if necessary
			amountInCents = int(payment.Amount * 100)
		}

		// Create the PaymentIntent with Stripe
		paymtID, clientSecret, err := m.StripePay.CreatePaymentIntent(float64(amountInCents), "usd")
		if err != nil {
			return nil, fmt.Errorf("failed to create payment intent: %v", err)
		}
		log.Printf("Payment Intent created with ID: %s", paymtID)

		// Update payment details
		payment.PaymentID = paymtID
		payment.ClientSecret = clientSecret
		payment.PaymentMethod = "stripe"
		payment.Status = "Pending" // Set status to "Pending" for the new payment intent

		// Update the order data in Redis
		updatedPaymentData, err := json.Marshal(payment)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal updated payment data: %v", err)
		}
		err = m.redis.SetDataInRedis(key, updatedPaymentData, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to update payment data in Redis: %v", err)
		}

		// Save the new payment intent details in the database
		err = m.Repo.SavePayment(&payment)
		if err != nil {
			return nil, fmt.Errorf("failed to save payment data in the database: %v", err)
		}
	}

	// Return the payment response using the existing or newly created PaymentIntent
	response := &pb.PaymentResponse{
		PaymentId:    payment.PaymentID,
		ClientSecret: payment.ClientSecret,
		OrderId:      fmt.Sprintf("%d", p.Order_ID),
		Amount:       float64(payment.Amount),
	}
	return response, nil
}
func (m *MaterialService) PaymentSuccessService(p *pb.Payment) (*pb.PaymentStatusResponse, error) {
	key := fmt.Sprintf("order:%d:user:%d", p.Order_ID, p.User_ID)
	log.Printf("Retrieving payment data from Redis for key: %s", key)

	// Fetch payment data from Redis
	paymentData, err := m.redis.GetFromRedis(key)
	var payment model.Payment
	if err == redis.Nil {
		log.Printf("No payment data / payment paid already; fetching from database for order_id: %d", p.Order_ID)

		// Fetch the latest payment from the database
		payment, err = m.Repo.GetLatestPaymentByOrderID(int(p.Order_ID))
		if err != nil {
			return &pb.PaymentStatusResponse{
				Status:  pb.PaymentStatusResponse_FAILED,
				Message: fmt.Sprintf("Failed to fetch payment data from the database: %v", err),
			}, nil
		}
	} else if err != nil {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: fmt.Sprintf("Error fetching payment data from Redis: %v", err),
		}, nil
	} else {
		json.Unmarshal([]byte(paymentData), &payment)
	}

	// Check if payment has already been completed
	if payment.Status == "Completed" {
		log.Printf("Payment already completed for payment_id: %s", payment.PaymentID)
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_SUCCESS,
			Message: "Payment has already been successfully processed.",
		}, nil
	}

	// Verify payment status with Stripe
	paymentStatus, err := m.StripePay.VerifyPaymentStatus(payment.PaymentID)
	if err != nil {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: fmt.Sprintf("Failed to verify payment status: %v", err),
		}, nil
	}

	// If payment is not successful, return failure
	if paymentStatus != "succeeded" {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: "Payment failed or not completed.",
		}, nil
	}

	// Check if the order is already marked as completed before updating
	order, err := m.Repo.FindOrdersByID(uint(p.Order_ID))
	if err != nil {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: fmt.Sprintf("Error fetching order data from the database: %v", err),
		}, nil
	}

	// Prevent updating the payment if the order is already completed
	if order.Status == "Completed" {
		log.Printf("Order already completed, skipping payment update for order_id: %d", p.Order_ID)
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_SUCCESS,
			Message: "Order has already been completed and paid.",
		}, nil
	}

	// Update payment and order status in the database
	payment.Status = "Completed"
	err = m.Repo.UpdatePaymentAndOrderStatus(payment.PaymentID, int(p.Order_ID), payment.Status, "Completed")
	if err != nil {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: fmt.Sprintf("Failed to update payment and order status: %v", err),
		}, nil
	}

	// Publish Payment Success Event
	err = utils.HandlePaymentNotification(payment.PaymentID, uint(p.Order_ID), order.Email, payment.Amount, time.Now())
	if err != nil {
		log.Printf("Failed to publish payment event: %v", err)
	}

	// Update Redis with the latest payment data
	updatedPaymentData, _ := json.Marshal(payment)
	err = m.redis.SetDataInRedis(key, updatedPaymentData, 0) // Ensure that the cache is updated
	if err != nil {
		log.Printf("Failed to update Redis cache: %v", err)
	}

	// remove the payment data from Redis if no longer needed
	err = m.redis.DeleteDataFromRedis(key)
	if err != nil {
		log.Printf("Failed to delete Redis cache after payment update: %v", err)
	}

	// Fetch the item for the order (assuming each order has only one item for cutting result generation)
	item, err := m.Repo.FindItemByID(order.ItemID)
	if err != nil {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: fmt.Sprintf("Error loading item for order ID %d: %v", p.Order_ID, err),
		}, nil
	}

	log.Printf("Triggering cutting result generation for item ID: %d", item.ID)

	// Generate cutting result for the item
	_, err = m.GenerateCuttingResult(&pb.ItemID{ID: uint32(item.ID)})
	if err != nil {
		return &pb.PaymentStatusResponse{
			Status:  pb.PaymentStatusResponse_FAILED,
			Message: fmt.Sprintf("Failed to generate cutting result for item ID %d: %v", item.ID, err),
		}, nil
	}

	// fmt.Println("cutting itemID", pb.ItemID{ID: uint32(order.ItemID)})
	// _, err = m.GenerateCuttingResult(&pb.ItemID{ID: uint32(order.ItemID)})
	// if err != nil {
	// 	return &pb.PaymentStatusResponse{
	// 		Status:  pb.PaymentStatusResponse_FAILED,
	// 		Message: fmt.Sprintf("Failed to generate cutting result: %v", err),
	// 	}, nil
	// }

	log.Printf("Payment successfully completed and updated for payment_id: %s", payment.PaymentID)
	return &pb.PaymentStatusResponse{
		Status:  pb.PaymentStatusResponse_SUCCESS,
		Message: "Payment successfully completed.",
	}, nil
}
