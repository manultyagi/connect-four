package main

const (
	Rows    = 6
	Columns = 7
)

// Game represents a Connect Four game
type Game struct {
	Board  [Rows][Columns]int // 0 = empty, 1 = player1, 2 = player2
	Turn   int                // current player (1 or 2)
	Winner int                // 0 = none, 1 or 2 = winner, -1 = draw
	Moves  int                // total moves made
}

// NewGame creates a new game
func NewGame() *Game {
	return &Game{
		Turn: 1,
	}
}

// DropDisc drops a disc into a column
func (g *Game) DropDisc(col int) bool {
	if col < 0 || col >= Columns || g.Winner != 0 {
		return false
	}

	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = g.Turn
			g.Moves++
			return true
		}
	}
	return false
}

// SwitchTurn switches the current player
func (g *Game) SwitchTurn() {
	if g.Turn == 1 {
		g.Turn = 2
	} else {
		g.Turn = 1
	}
}

// Horizontal win check
func (g *Game) checkHorizontal(player int) bool {
	for row := 0; row < Rows; row++ {
		count := 0
		for col := 0; col < Columns; col++ {
			if g.Board[row][col] == player {
				count++
				if count == 4 {
					return true
				}
			} else {
				count = 0
			}
		}
	}
	return false
}

// Vertical win check
func (g *Game) checkVertical(player int) bool {
	for col := 0; col < Columns; col++ {
		count := 0
		for row := 0; row < Rows; row++ {
			if g.Board[row][col] == player {
				count++
				if count == 4 {
					return true
				}
			} else {
				count = 0
			}
		}
	}
	return false
}

// Diagonal ↘ (top-left to bottom-right)
func (g *Game) checkDiagonalDown(player int) bool {
	for row := 0; row <= Rows-4; row++ {
		for col := 0; col <= Columns-4; col++ {
			if g.Board[row][col] == player &&
				g.Board[row+1][col+1] == player &&
				g.Board[row+2][col+2] == player &&
				g.Board[row+3][col+3] == player {
				return true
			}
		}
	}
	return false
}

// Diagonal ↗ (bottom-left to top-right)
func (g *Game) checkDiagonalUp(player int) bool {
	for row := 3; row < Rows; row++ {
		for col := 0; col <= Columns-4; col++ {
			if g.Board[row][col] == player &&
				g.Board[row-1][col+1] == player &&
				g.Board[row-2][col+2] == player &&
				g.Board[row-3][col+3] == player {
				return true
			}
		}
	}
	return false
}

// CheckWin checks all win conditions
func (g *Game) CheckWin(player int) bool {
	return g.checkHorizontal(player) ||
		g.checkVertical(player) ||
		g.checkDiagonalDown(player) ||
		g.checkDiagonalUp(player)
}

// IsDraw checks if the game is a draw
func (g *Game) IsDraw() bool {
	return g.Moves == Rows*Columns && g.Winner == 0
}

// MakeMove performs a full move cycle
func (g *Game) MakeMove(col int) bool {
	if !g.DropDisc(col) {
		return false
	}

	if g.CheckWin(g.Turn) {
		g.Winner = g.Turn
	} else if g.IsDraw() {
		g.Winner = -1
	} else {
		g.SwitchTurn()
	}

	return true
}
