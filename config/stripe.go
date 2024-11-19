// package config

// import (
// 	"fmt"
// 	"log"

// 	"github.com/stripe/stripe-go/v72"
// 	"github.com/stripe/stripe-go/v72/paymentintent"
// )

// type StripeClient struct {
// 	apiKey string
// }

// // NewStripeClient initializes the Stripe client with the API key from the configuration.
// func NewStripeClient(cfg Config) *StripeClient {
// 	stripe.Key = cfg.STRIPESECRETKEY
// 	if stripe.Key == "" {
// 		log.Fatal("Stripe secret key is missing from configuration")
// 	}
// 	//log.Println("my stripe key", stripe.Key)
// 	return &StripeClient{
// 		apiKey: cfg.STRIPESECRETKEY,
// 	}
// }

// // CreatePaymentIntent creates a Stripe Payment Intent and returns its ID.
// func (s *StripeClient) CreatePaymentIntent(amount float64, currency string) (string, error) {
// 	// Convert amount to smallest currency unit
// 	fmt.Println("amount and currency", amount, currency)
// 	amountInSmallestUnit, err := convertToSmallestUnit(amount, currency)
// 	if err != nil {
// 		return "", fmt.Errorf("invalid amount or currency: %v", err)
// 	}

// 	// Validate the minimum amount
// 	if !isValidAmount(amountInSmallestUnit, currency) {
// 		return "", fmt.Errorf("amount is below the minimum allowed amount for the currency")
// 	}

// 	params := &stripe.PaymentIntentParams{
// 		Amount:   stripe.Int64(amountInSmallestUnit),
// 		Currency: stripe.String(currency),
// 		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
// 			Enabled: stripe.Bool(true),
// 		},
// 	}

// 	// Create the PaymentIntent
// 	pi, err := paymentintent.New(params)
// 	if err != nil {
// 		log.Printf("Stripe API error: %v", err)
// 		return "", fmt.Errorf("failed to create payment intent: %v", err)
// 	}

// 	// params := &stripe.PaymentIntentParams{
// 	// 	Amount:   stripe.Int64(amountInSmallestUnit),
// 	// 	Currency: stripe.String(currency),
// 	// 	AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
// 	// 		Enabled: stripe.Bool(true),
// 	// 	},
// 	// }

// 	// // Create the PaymentIntent
// 	// pi, err := paymentintent.New(params)
// 	// if err != nil {
// 	// 	if stripeErr, ok := err.(*stripe.Error); ok {
// 	// 		fmt.Printf("A stripe error occured:%v\n", stripeErr.Error())
// 	// 	}
// 	// 	log.Printf("Stripe API error: %v", err)
// 	// 	return "", fmt.Errorf("failed to create payment intent: %v", err)
// 	// }

// 	return pi.ClientSecret, nil
// }

// // convertToSmallestUnit converts the amount into the smallest unit based on the currency (e.g., cents for USD).
// func convertToSmallestUnit(amount float64, currency string) (int64, error) {
// 	switch currency {
// 	case "usd":
// 		// For USD, convert amount to cents
// 		return int64(amount * 100), nil
// 	case "inr":
// 		// For INR, convert amount to paise (100 paise = 1 INR)
// 		return int64(amount * 100), nil
// 	default:
// 		// If the currency is not supported or unknown, return an error
// 		return 0, fmt.Errorf("unsupported currency: %s", currency)
// 	}
// }

// // isValidAmount checks if the amount is above the minimum allowed for the currency.
// func isValidAmount(amount int64, currency string) bool {
// 	// Define minimum amounts for each currency
// 	switch currency {
// 	case "usd":
// 		// Minimum for USD is $0.50 (i.e., 50 cents)
// 		return amount >= 50
// 	case "inr":
// 		// Minimum for INR is ₹50 (i.e., 5000 paise)
// 		return amount >= 5000
// 	default:
// 		return false
// 	}
// }

// // VerifyPaymentStatus checks the status of a payment intent.
// func (s *StripeClient) VerifyPaymentStatus(paymentID string) (string, error) {
// 	stripe.Key = s.apiKey // Use StripeClient’s apiKey
// 	log.Print("paymentID", paymentID)
// 	intent, err := paymentintent.Get(paymentID, nil)
// 	if err != nil {
// 		return "", fmt.Errorf("error retrieving payment intent: %v", err)
// 	}

//		return string(intent.Status), nil
//	}
package config

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

// StripeClient holds the Stripe configuration and methods to interact with the Stripe API.
type StripeClient struct {
	apiKey string
	redis  *RedisService // Assuming a Redis client is used
}

// NewStripeClient initializes the Stripe client with the API key from the configuration.
func NewStripeClient(cfg Config, redis *RedisService) *StripeClient {
	stripe.Key = cfg.STRIPESECRETKEY
	if stripe.Key == "" {
		log.Fatal("Stripe secret key is missing from configuration")
	}
	return &StripeClient{
		apiKey: cfg.STRIPESECRETKEY,
		redis:  redis,
	}
}

// CreatePaymentIntent creates a Stripe PaymentIntent and returns its ID and client secret.
func (s *StripeClient) CreatePaymentIntent(amount float64, currency string) (string, string, error) {
	// Convert amount to the smallest unit (e.g., cents for USD)
	amountInSmallestUnit, err := convertToSmallestUnit(amount, currency)
	if err != nil {
		return "", "", fmt.Errorf("invalid amount or currency: %v", err)
	}

	// Validate the minimum allowed amount for the given currency
	if !isValidAmount(amountInSmallestUnit, currency) {
		return "", "", fmt.Errorf("amount is below the minimum allowed amount for the currency")
	}

	// Create the PaymentIntent with the specified amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInSmallestUnit),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Create the PaymentIntent
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Stripe API error: %v", err)
		return "", "", fmt.Errorf("failed to create payment intent: %v", err)
	}

	// Store the payment ID and client secret in Redis for later use
	err = s.storeClientSecretInRedis(pi.ID, pi.ClientSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to store client secret in Redis: %v", err)
	}

	// Return the PaymentIntent ID and client secret
	return pi.ID, pi.ClientSecret, nil
}

// storeClientSecretInRedis stores the PaymentIntent's client secret in Redis.
func (s *StripeClient) storeClientSecretInRedis(paymentID, clientSecret string) error {
	//err := s.redis.SetDataInRedis(fmt.Sprintf("payment:%s", paymentID), clientSecret, 0)
	err := s.redis.SetDataInRedis(fmt.Sprintf("payment:%s", paymentID), []byte(clientSecret), 0)

	if err != nil {
		log.Printf("Failed to store client secret for payment %s in Redis: %v", paymentID, err)
		return err
	}
	return nil
}

// VerifyPaymentStatus retrieves the status of a payment intent by its ID.
func (s *StripeClient) VerifyPaymentStatus(paymentID string) (string, error) {
	stripe.Key = s.apiKey // Use the client's Stripe API key

	intent, err := paymentintent.Get(paymentID, nil)
	if err != nil {
		log.Printf("Failed to retrieve payment intent: %v", err)
		return "", fmt.Errorf("error retrieving payment intent: %v", err)
	}

	// Return the status of the payment intent (e.g., succeeded, failed, etc.)
	return string(intent.Status), nil
}

// convertToSmallestUnit converts the amount into the smallest unit based on the currency (e.g., cents for USD).
func convertToSmallestUnit(amount float64, currency string) (int64, error) {
	switch currency {
	case "usd":
		// For USD, convert amount to cents (100 cents = 1 dollar)
		return int64(amount * 100), nil
	case "inr":
		// For INR, convert amount to paise (100 paise = 1 INR)
		return int64(amount * 100), nil
	default:
		// If the currency is unsupported, return an error
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}
}

// isValidAmount checks if the amount is above the minimum allowed for the given currency.
func isValidAmount(amount int64, currency string) bool {
	switch currency {
	case "usd":
		// Minimum for USD is $0.50 (i.e., 50 cents)
		return amount >= 50
	case "inr":
		// Minimum for INR is ₹50 (i.e., 5000 paise)
		return amount >= 5000
	default:
		// For unsupported currencies, return false
		return false
	}
}
