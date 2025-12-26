package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func initKafka() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "game-events",
		Balancer: &kafka.LeastBytes{},
	}
}

func emitEvent(eventType string, payload interface{}) {
	event := map[string]interface{}{
		"type":      eventType,
		"payload":   payload,
		"timestamp": time.Now().UTC(),
	}

	bytes, _ := json.Marshal(event)

	_ = kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{Value: bytes},
	)
}
