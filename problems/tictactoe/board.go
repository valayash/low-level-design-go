package tictactoe

import (
	"fmt"
	"strings"
)

// Board stores the grid and nothing else. It knows how cells are laid out and
// how many in a row constitute a line, but it has no idea that players exist,
// whose turn it is, or whether a game is over.
type Board struct {
	grid      [][]Symbol
	winLength int
}

// NewBoard returns a ready-to-use size x size board. There is no separate
// initialize step, so a caller cannot obtain a half-built board.
func NewBoard(size, winLength int) (*Board, error) {
	if size < 1 {
		return nil, ErrInvalidSize
	}
	if winLength < 1 || winLength > size {
		return nil, ErrInvalidWinLength
	}

	grid := make([][]Symbol, size)
	for row := range grid {
		grid[row] = make([]Symbol, size)
		for col := range grid[row] {
			grid[row][col] = Empty
		}
	}

	return &Board{grid: grid, winLength: winLength}, nil
}

// Size is derived from the grid rather than stored, so it cannot drift.
func (b *Board) Size() int { return len(b.grid) }

// WinLength is how many consecutive marks win a line.
func (b *Board) WinLength() int { return b.winLength }

func (b *Board) inBounds(row, col int) bool {
	n := len(b.grid)
	return row >= 0 && row < n && col >= 0 && col < n
}

// GetCell is the only way to read a cell from outside the package, which is
// what lets the win rules live somewhere other than Board.
func (b *Board) GetCell(row, col int) (Symbol, error) {
	if !b.inBounds(row, col) {
		return Empty, fmt.Errorf("cell (%d,%d): %w", row, col, ErrOutOfBounds)
	}
	return b.grid[row][col], nil
}

// Mark places a symbol, rejecting anything that would corrupt the grid. It
// validates before it writes, so a rejected move leaves the board untouched.
func (b *Board) Mark(row, col int, s Symbol) error {
	if !b.inBounds(row, col) {
		return fmt.Errorf("cell (%d,%d): %w", row, col, ErrOutOfBounds)
	}
	if b.grid[row][col] != Empty {
		return fmt.Errorf("cell (%d,%d): %w", row, col, ErrCellOccupied)
	}

	b.grid[row][col] = s
	return nil
}

// IsFull reports whether any cell is still Empty. Scanning costs O(n^2) but
// cannot go stale the way a stored move counter can.
func (b *Board) IsFull() bool {
	for _, row := range b.grid {
		for _, cell := range row {
			if cell == Empty {
				return false
			}
		}
	}
	return true
}

// countDir walks outward from (row,col) along (dr,dc) and counts how many
// consecutive cells hold s. The starting cell itself is not counted.
func (b *Board) countDir(row, col, dr, dc int, s Symbol) int {
	count := 0
	r, c := row+dr, col+dc
	for b.inBounds(r, c) && b.grid[r][c] == s {
		count++
		r, c = r+dr, c+dc
	}
	return count
}

// IsWinningMove reports whether the mark at (row,col) completes a line of
// winLength. Only lines through that cell can have changed, so this checks
// four axes rather than rescanning the whole board.
func (b *Board) IsWinningMove(row, col int) bool {
	if !b.inBounds(row, col) {
		return false
	}

	s := b.grid[row][col]
	if s == Empty {
		return false
	}

	// horizontal, vertical, and the two diagonals
	directions := [4][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range directions {
		total := 1 +
			b.countDir(row, col, d[0], d[1], s) +
			b.countDir(row, col, -d[0], -d[1], s)
		if total >= b.winLength {
			return true
		}
	}

	return false
}

// String renders the grid. Formatting lives here so the rules never import
// anything that writes to a terminal.
func (b *Board) String() string {
	var sb strings.Builder
	for _, row := range b.grid {
		for col, cell := range row {
			if col > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteRune(rune(cell))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
