package config

import (
	"fmt"

	razorpay "github.com/razorpay/razorpay-go"
)

// RazorpayClient wraps the Razorpay SDK client.
type RazorpayClient struct {
	client *razorpay.Client
}

// // Config should contain your API credentials for Razorpay.
// type Config struct {
// 	APIKey    string
// 	APISecret string
// }

// NewRazorpayClient initializes and returns a Razorpay client.
func NewRazorpayClient(cfg Config) *RazorpayClient {
	client := razorpay.NewClient(cfg.APIKey, cfg.APISecret)
	return &RazorpayClient{
		client: client,
	}
}

// CreateOrder generates a new Razorpay order with the specified amount (in INR).
func (r *RazorpayClient) CreateOrder(amount float64) (string, error) {
	// Convert the amount to paisa (1 INR = 100 paisa)
	amountInPaisa := int(amount * 100)

	// Prepare the order data with required fields
	data := map[string]interface{}{
		"amount":          amountInPaisa,
		"currency":        "INR",
		"receipt":         "opti_cutter", // Customize the receipt number as needed
		"payment_capture": 1,             // Auto-capture the payment once created
	}

	// Create order using Razorpay API
	body, err := r.client.Order.Create(data, nil)
	if err != nil {
		return "", fmt.Errorf("problem creating order: %v", err)
	}

	// Extract and return the order ID
	// orderID, ok := body["id"].(string)
	// if !ok {
	// 	return "", fmt.Errorf("unable to retrieve order ID from response")
	// }
	value := body["id"]
	orderID := value.(string)

	return orderID, nil
}
