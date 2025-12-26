package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// --------------------
// Models
// --------------------

type Player struct {
	Username     string
	Conn         *websocket.Conn
	PlayerNumber int
}

type ClientMessage struct {
	Column int `json:"column"`
}

type ServerMessage struct {
	Board        [Rows][Columns]int `json:"board"`
	Turn         int                `json:"turn"`
	Winner       int                `json:"winner"`
	PlayerNumber int                `json:"playerNumber"`
}

type GameSession struct {
	Game    *Game
	Player1 *Player
	Player2 *Player // nil = bot
}

type PlayerMove struct {
	Player *Player
	Column int
}

// --------------------
// Globals
// --------------------

var waitingPlayer *Player

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// --------------------
// WebSocket Handler
// --------------------

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	var intro struct {
		Username string `json:"username"`
	}

	if err := conn.ReadJSON(&intro); err != nil {
		conn.Close()
		return
	}

	player := &Player{
		Username: intro.Username,
		Conn:     conn,
	}

	if waitingPlayer == nil {
		player.PlayerNumber = 1
		waitingPlayer = player

		conn.WriteJSON(map[string]string{
			"status": "waiting for opponent",
		})

		go func() {
			time.Sleep(10 * time.Second)
			if waitingPlayer == player {
				waitingPlayer = nil
				runGameSession(&GameSession{
					Game:    NewGame(),
					Player1: player,
					Player2: nil,
				})
			}
		}()

	} else {
		player.PlayerNumber = 2
		opponent := waitingPlayer
		waitingPlayer = nil

		runGameSession(&GameSession{
			Game:    NewGame(),
			Player1: opponent,
			Player2: player,
		})
	}
}

// --------------------
// Game Session
// --------------------

func runGameSession(session *GameSession) {
	game := session.Game
	moves := make(chan PlayerMove)

	player2Name := "BOT"
	if session.Player2 != nil {
		player2Name = session.Player2.Username
	}

	gameID := fmt.Sprintf("%d-%s", time.Now().UnixNano(), session.Player1.Username)
	startTime := time.Now().UTC()

	emitEvent("GAME_STARTED", map[string]interface{}{
		"gameId":  gameID,
		"player1": session.Player1.Username,
		"player2": player2Name,
	})

	sendState := func() {
		send := func(p *Player) {
			if p == nil {
				return
			}
			p.Conn.WriteJSON(ServerMessage{
				Board:        game.Board,
				Turn:         game.Turn,
				Winner:       game.Winner,
				PlayerNumber: p.PlayerNumber,
			})
		}
		send(session.Player1)
		send(session.Player2)
	}

	startReader := func(p *Player) {
		if p == nil {
			return
		}
		go func() {
			for {
				var msg ClientMessage
				if err := p.Conn.ReadJSON(&msg); err != nil {
					return
				}
				moves <- PlayerMove{Player: p, Column: msg.Column}
			}
		}()
	}

	startReader(session.Player1)
	startReader(session.Player2)

	sendState()

	for game.Winner == 0 {

		// 🤖 Bot move
		if game.Turn == 2 && session.Player2 == nil {
			time.Sleep(800 * time.Millisecond)
			game.MakeMove(BotMove(game))
			sendState()
			continue
		}

		move := <-moves

		// 🔐 Enforce turn
		if move.Player.PlayerNumber != game.Turn {
			continue
		}

		if game.MakeMove(move.Column) {
			sendState()
		}
	}

	saveGameResult(session)

	emitEvent("GAME_FINISHED", map[string]interface{}{
		"gameId":    gameID,
		"player1":   session.Player1.Username,
		"player2":   player2Name,
		"winner":    game.Winner,
		"moves":     game.Moves,
		"startedAt": startTime,
		"endedAt":   time.Now().UTC(),
	})

	fmt.Println("Game finished")
}

// --------------------
// Leaderboard
// --------------------

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := getLeaderboard()
	if err != nil {
		http.Error(w, "failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}
