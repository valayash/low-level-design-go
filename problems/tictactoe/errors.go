package tictactoe

import "errors"

// Sentinel errors let callers distinguish *why* an operation failed using
// errors.Is, instead of collapsing every failure into a single false.
var (
	ErrOutOfBounds      = errors.New("position is outside the board")
	ErrCellOccupied     = errors.New("cell is already occupied")
	ErrGameOver         = errors.New("game is already over")
	ErrInvalidSize      = errors.New("board size must be at least 1")
	ErrInvalidWinLength = errors.New("win length must be between 1 and the board size")
	ErrTooFewPlayers    = errors.New("a game needs at least two players")
	ErrNilPlayer        = errors.New("player must not be nil")
	ErrEmptySymbol      = errors.New("a player's sign must not be Empty")
	ErrDuplicateSymbol  = errors.New("players must have distinct signs")
)
