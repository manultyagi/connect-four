package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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
	Status       string             `json:"status,omitempty"`
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

var waitingPlayer *Player

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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
		Username:     intro.Username,
		Conn:         conn,
		PlayerNumber: 1,
	}

	if waitingPlayer == nil {
		waitingPlayer = player

		player.Conn.WriteJSON(ServerMessage{
			Status: "waiting for opponent",
		})

		go func(p *Player) {
			time.Sleep(10 * time.Second)
			if waitingPlayer == p {
				waitingPlayer = nil
				runGameSession(&GameSession{
					Game:    NewGame(),
					Player1: p,
					Player2: nil,
				})
			}
		}(player)

	} else {
		opponent := waitingPlayer
		waitingPlayer = nil

		player.PlayerNumber = 2
		opponent.PlayerNumber = 1

		runGameSession(&GameSession{
			Game:    NewGame(),
			Player1: opponent,
			Player2: player,
		})
	}
}

func runGameSession(session *GameSession) {
	game := session.Game
	moves := make(chan PlayerMove)

	sendState := func(status string) {
		send := func(p *Player) {
			if p == nil {
				return
			}
			p.Conn.WriteJSON(ServerMessage{
				Board:        game.Board,
				Turn:         game.Turn,
				Winner:       game.Winner,
				PlayerNumber: p.PlayerNumber,
				Status:       status,
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

	// 🔑 EXPLICIT GAME START STATE
	sendState("game started")

	for game.Winner == 0 {

		// Bot move
		if game.Turn == 2 && session.Player2 == nil {
			time.Sleep(700 * time.Millisecond)
			game.MakeMove(BotMove(game))
			sendState("")
			continue
		}

		move := <-moves

		if move.Player.PlayerNumber != game.Turn {
			continue
		}

		if game.MakeMove(move.Column) {
			sendState("")
		}
	}

	sendState("game finished")
	saveGameResult(session)
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := getLeaderboard()
	if err != nil {
		http.Error(w, "failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}
