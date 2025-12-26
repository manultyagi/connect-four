package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Global DB connection
var db *sql.DB

// InitDB initializes the PostgreSQL connection
func InitDB() {
	var err error

	connStr := "user=postgres password=Iamnulagi7 dbname=connect_four sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("Database connected")
}

// saveGameResult persists a completed game and updates player wins
func saveGameResult(session *GameSession) {
	// 🔑 CRITICAL FIX:
	// If DB is disabled (db == nil), do nothing
	if db == nil {
		log.Println("DB disabled, skipping saveGameResult")
		return
	}

	game := session.Game

	var winner string

	if game.Winner == 1 {
		winner = session.Player1.Username
	} else if game.Winner == 2 {
		if session.Player2 != nil {
			winner = session.Player2.Username
		} else {
			winner = "BOT"
		}
	} else {
		winner = "DRAW"
	}

	player2 := "BOT"
	if session.Player2 != nil {
		player2 = session.Player2.Username
	}

	// Save game record
	_, err := db.Exec(
		`INSERT INTO games (player1, player2, winner, moves)
		 VALUES ($1, $2, $3, $4)`,
		session.Player1.Username,
		player2,
		winner,
		game.Moves,
	)
	if err != nil {
		log.Println("Failed to save game:", err)
		return
	}

	// Update player wins (human only)
	if winner != "DRAW" && winner != "BOT" {
		_, err = db.Exec(
			`INSERT INTO players (username, wins)
			 VALUES ($1, 1)
			 ON CONFLICT (username)
			 DO UPDATE SET wins = players.wins + 1`,
			winner,
		)
		if err != nil {
			log.Println("Failed to update player wins:", err)
		}
	}
}
