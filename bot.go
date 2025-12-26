package main

//simulate check if dropping a disc causes a win
func simulateMove(board [Rows][Columns]int, col int, player int) bool {
	for row := Rows - 1; row >= 0; row-- {
		if board[row][col] == 0 {
			board[row][col] = player
			tempGame := &Game{Board: board}
			return tempGame.CheckWin((player))
		}
	}
	return false
}

// BotMove decides the best column for the bot to play
// Bot is always player 2
func BotMove(g *Game) int {

	// 1️⃣ Try to win
	for col := 0; col < Columns; col++ {
		if simulateMove(g.Board, col, 2) {
			return col
		}
	}

	// 2️⃣ Block opponent (player 1)
	for col := 0; col < Columns; col++ {
		if simulateMove(g.Board, col, 1) {
			return col
		}
	}

	// 3️⃣ Prefer center column
	center := Columns / 2
	if g.Board[0][center] == 0 {
		return center
	}

	// 4️⃣ Fallback: first available column
	for col := 0; col < Columns; col++ {
		if g.Board[0][col] == 0 {
			return col
		}
	}

	return -1 // no valid moves
}
