package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
}

var (
	totalGames    int
	totalDuration time.Duration
	winsPerPlayer = map[string]int{}
	gamesPerHour  = map[string]int{}
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "game-events",
		GroupID: "analytics-service",
	})

	fmt.Println("Analytics service started...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}

		handleEvent(event)
	}
}

func handleEvent(event Event) {
	switch event.Type {

	case "GAME_FINISHED":
		totalGames++

		startedAt, _ := time.Parse(time.RFC3339, fmt.Sprint(event.Payload["startedAt"]))
		endedAt, _ := time.Parse(time.RFC3339, fmt.Sprint(event.Payload["endedAt"]))

		duration := endedAt.Sub(startedAt)
		totalDuration += duration

		winner := fmt.Sprint(event.Payload["winner"])
		if winner != "-1" {
			winsPerPlayer[winner]++
		}

		hourKey := endedAt.Format("2006-01-02 15")
		gamesPerHour[hourKey]++

		fmt.Println("---- Analytics ----")
		fmt.Println("Total games:", totalGames)
		fmt.Println("Avg game duration:", totalDuration/time.Duration(totalGames))
		fmt.Println("Wins per player:", winsPerPlayer)
		fmt.Println("Games per hour:", gamesPerHour)
		fmt.Println("-------------------")
	}
}
