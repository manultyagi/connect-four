package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// 1️⃣ Initialize database
	if os.Getenv("ENABLE_DB") == "true" {
		fmt.Println("Database enabled")
		InitDB()
	} else {
		fmt.Println("Database disabled")
	}

	// 2️⃣ Initialize Kafka ONLY if enabled
	if os.Getenv("ENABLE_KAFKA") == "true" {
		fmt.Println("Kafka enabled")
		initKafka()
	} else {
		fmt.Println("Kafka disabled")
	}

	// 3️⃣ WebSocket + HTTP routes
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/leaderboard", leaderboardHandler)

	// 4️⃣ Serve frontend
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// 5️⃣ Dynamic port (required for hosting)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started on port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
