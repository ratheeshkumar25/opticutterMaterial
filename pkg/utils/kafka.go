package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	OrderID   uint    `json:"order_id"`
	UserID    uint    `json:"user_id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
}

func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        "payment_topic",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireOne),
	})
	return &KafkaProducer{writer: writer}, nil
}

func (k *KafkaProducer) ProducerPaymentEvent(event PaymentEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event:%w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event.PaymentID),
		Value: message,
	}
	//send the message

	err = k.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	log.Printf("Message successfully sent to Kafka topic: %s", k.writer.Topic)
	return nil
}

func HandlePaymentNotification(paymentID string, orderID uint, email string, amount float64, datetime time.Time) error {
	kafkaProducer, err := NewKafkaProducer("kafka:9092")
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Create a payment event
	paymentEvent := PaymentEvent{
		PaymentID: paymentID,
		OrderID:   orderID,
		Email:     email,
		Amount:    amount,
		Date:      time.Now().Format(time.RFC3339),
	}

	fmt.Println("Payment event is topic is", paymentEvent)
	err = kafkaProducer.ProducerPaymentEvent(paymentEvent)
	if err != nil {
		return fmt.Errorf("failed to produce payment event: %w", err)
	}

	log.Println("Payment event produced successfully")
	return nil
}
