package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type CuttingResultEvent struct {
	CuttingResultID uint               `json:"cutting_result_id"`
	ItemID          uint               `json:"item_id"`
	Components      []ComponentPayload `json:"components"`
}

type ComponentPayload struct {
	MaterialID    uint   `json:"material_id"`
	DoorPanel     string `json:"door_panel,omitempty"`
	BackSidePanel string `json:"back_side_panel,omitempty"`
	SidePanel     string `json:"side_panel,omitempty"`
	TopPanel      string `json:"top_panel,omitempty"`
	BottomPanel   string `json:"bottom_panel,omitempty"`
	ShelvesPanel  string `json:"shelves_panel,omitempty"`
	PanelCount    int32  `json:"panel_count"`
}

type KafkaCuttingResultProducer struct {
	writer *kafka.Writer
}

func NewKafkaCuttingResultProducer(broker string) (*KafkaCuttingResultProducer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        "cutting_topic",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireOne),
	})
	return &KafkaCuttingResultProducer{writer: writer}, nil
}

func (k *KafkaCuttingResultProducer) ProducerCuttingResultEvent(event CuttingResultEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", event.CuttingResultID)),
		Value: message,
	}
	err = k.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to produce cutting result event: %w", err)
	}
	log.Printf("Cutting result event successfully sent to Kafka topic: %s", k.writer.Topic)
	return nil
}

func HandleCuttingResultNotification(cuttingResultID uint, components []ComponentPayload, itemID uint) error {
	kafkaProducer, err := NewKafkaCuttingResultProducer("localhost:9092")
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Create and send CuttingResultEvent
	cuttingResultEvent := CuttingResultEvent{
		CuttingResultID: cuttingResultID,
		ItemID:          itemID,
		Components:      components,
	}
	err = kafkaProducer.ProducerCuttingResultEvent(cuttingResultEvent)
	if err != nil {
		return fmt.Errorf("failed to produce cutting result event: %w", err)
	}

	log.Println("Cutting result event produced successfully")
	return nil
}
