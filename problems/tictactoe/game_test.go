package tictactoe

import (
	"errors"
	"testing"
)

// newTestGame builds a standard 3x3, two-player game and fails the test if
// construction is rejected.
func newTestGame(t *testing.T) (*Game, *Player, *Player) {
	t.Helper()

	x := NewPlayer("Xavier", X)
	o := NewPlayer("Olivia", O)

	g, err := NewGame(3, 3, x, o)
	if err != nil {
		t.Fatalf("NewGame: unexpected error: %v", err)
	}
	return g, x, o
}

// playAll applies moves in order and fails on the first rejection.
func playAll(t *testing.T, g *Game, moves [][2]int) {
	t.Helper()

	for i, m := range moves {
		if err := g.Move(m[0], m[1]); err != nil {
			t.Fatalf("move %d at (%d,%d): unexpected error: %v", i, m[0], m[1], err)
		}
	}
}

func TestWinAcrossRow(t *testing.T) {
	g, x, _ := newTestGame(t)

	// X takes the top row; O answers in the middle row.
	playAll(t, g, [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}})

	if g.Status() != Won {
		t.Fatalf("status = %v, want %v", g.Status(), Won)
	}
	w, ok := g.Winner()
	if !ok || w != x {
		t.Fatalf("winner = %v (ok=%v), want %v", w, ok, x)
	}
}

func TestWinDownColumn(t *testing.T) {
	g, x, _ := newTestGame(t)

	playAll(t, g, [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}, {2, 0}})

	if g.Status() != Won {
		t.Fatalf("status = %v, want %v", g.Status(), Won)
	}
	if w, _ := g.Winner(); w != x {
		t.Fatalf("winner = %v, want %v", w, x)
	}
}

func TestWinOnMainDiagonal(t *testing.T) {
	g, x, _ := newTestGame(t)

	playAll(t, g, [][2]int{{0, 0}, {0, 1}, {1, 1}, {0, 2}, {2, 2}})

	if g.Status() != Won {
		t.Fatalf("status = %v, want %v", g.Status(), Won)
	}
	if w, _ := g.Winner(); w != x {
		t.Fatalf("winner = %v, want %v", w, x)
	}
}

func TestWinOnAntiDiagonal(t *testing.T) {
	g, x, _ := newTestGame(t)

	playAll(t, g, [][2]int{{0, 2}, {0, 0}, {1, 1}, {0, 1}, {2, 0}})

	if g.Status() != Won {
		t.Fatalf("status = %v, want %v", g.Status(), Won)
	}
	if w, _ := g.Winner(); w != x {
		t.Fatalf("winner = %v, want %v", w, x)
	}
}

// A full board with no line is a draw.
//
//	X O X
//	X O O
//	O X X
func TestDrawOnFullBoard(t *testing.T) {
	g, _, _ := newTestGame(t)

	playAll(t, g, [][2]int{
		{0, 0}, {0, 1}, {0, 2}, {1, 1}, {1, 0},
		{1, 2}, {2, 2}, {2, 0}, {2, 1},
	})

	if g.Status() != Draw {
		t.Fatalf("status = %v, want %v\n%s", g.Status(), Draw, g.Board())
	}
	if _, ok := g.Winner(); ok {
		t.Fatal("a drawn game must have no winner")
	}
}

// The ordering trap: the final mark fills the board *and* completes a line.
// Checking for a draw before a win would report the wrong result.
//
//	O O X
//	X O O
//	X X X   <- X completes the bottom row with the last empty cell
func TestWinningMoveThatAlsoFillsBoard(t *testing.T) {
	g, x, _ := newTestGame(t)

	playAll(t, g, [][2]int{
		{0, 2}, {0, 0}, {1, 0}, {0, 1}, {2, 0},
		{1, 1}, {2, 1}, {1, 2}, {2, 2},
	})

	if !g.Board().IsFull() {
		t.Fatalf("board should be full:\n%s", g.Board())
	}
	if g.Status() != Won {
		t.Fatalf("status = %v, want %v — a winning final move is not a draw", g.Status(), Won)
	}
	if w, _ := g.Winner(); w != x {
		t.Fatalf("winner = %v, want %v", w, x)
	}
}

func TestOccupiedCellIsRejected(t *testing.T) {
	g, _, _ := newTestGame(t)

	playAll(t, g, [][2]int{{1, 1}})

	err := g.Move(1, 1)
	if !errors.Is(err, ErrCellOccupied) {
		t.Fatalf("err = %v, want ErrCellOccupied", err)
	}
}

func TestOutOfBoundsIsRejected(t *testing.T) {
	g, _, _ := newTestGame(t)

	// 3 is the classic off-by-one: valid indices are 0, 1, 2.
	for _, m := range [][2]int{{3, 0}, {0, 3}, {-1, 0}, {0, -1}} {
		if err := g.Move(m[0], m[1]); !errors.Is(err, ErrOutOfBounds) {
			t.Errorf("Move(%d,%d) err = %v, want ErrOutOfBounds", m[0], m[1], err)
		}
	}
}

// A rejected move must not consume a turn.
func TestRejectedMoveKeepsSameplayer(t *testing.T) {
	g, x, o := newTestGame(t)

	playAll(t, g, [][2]int{{0, 0}}) // X moves, turn passes to O

	if g.CurrentPlayer() != o {
		t.Fatalf("current = %v, want %v", g.CurrentPlayer(), o)
	}

	if err := g.Move(0, 0); err == nil { // O tries an occupied cell
		t.Fatal("expected the occupied cell to be rejected")
	}

	if g.CurrentPlayer() != o {
		t.Fatalf("current = %v, want %v — a rejected move must not advance the turn", g.CurrentPlayer(), o)
	}
	_ = x
}

func TestMoveAfterGameOverIsRejected(t *testing.T) {
	g, _, _ := newTestGame(t)

	playAll(t, g, [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}}) // X wins

	if err := g.Move(2, 2); !errors.Is(err, ErrGameOver) {
		t.Fatalf("err = %v, want ErrGameOver", err)
	}
}

func TestTurnsAlternate(t *testing.T) {
	g, x, o := newTestGame(t)

	if g.CurrentPlayer() != x {
		t.Fatalf("first mover = %v, want %v", g.CurrentPlayer(), x)
	}
	playAll(t, g, [][2]int{{0, 0}})
	if g.CurrentPlayer() != o {
		t.Fatalf("second mover = %v, want %v", g.CurrentPlayer(), o)
	}
	playAll(t, g, [][2]int{{1, 1}})
	if g.CurrentPlayer() != x {
		t.Fatalf("third mover = %v, want %v", g.CurrentPlayer(), x)
	}
}

// The generalisation the follow-ups ask for: a larger board where fewer than
// Size marks in a row wins.
func TestThreeInARowOnFourByFour(t *testing.T) {
	x := NewPlayer("Xavier", X)
	o := NewPlayer("Olivia", O)

	g, err := NewGame(4, 3, x, o)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}

	// X builds (1,1),(1,2),(1,3) — three in a row on a four-wide board.
	playAll(t, g, [][2]int{{1, 1}, {0, 0}, {1, 2}, {0, 1}, {1, 3}})

	if g.Status() != Won {
		t.Fatalf("status = %v, want %v\n%s", g.Status(), Won, g.Board())
	}
	if w, _ := g.Winner(); w != x {
		t.Fatalf("winner = %v, want %v", w, x)
	}
}

// Adding a third player costs nothing because turn order is a modulo.
func TestThreePlayersRotate(t *testing.T) {
	a := NewPlayer("A", X)
	b := NewPlayer("B", O)
	c := NewPlayer("C", Symbol('Z'))

	g, err := NewGame(5, 3, a, b, c)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}

	want := []*Player{a, b, c, a}
	moves := [][2]int{{0, 0}, {0, 1}, {0, 2}, {4, 4}}

	for i, m := range moves {
		if g.CurrentPlayer() != want[i] {
			t.Fatalf("move %d: current = %v, want %v", i, g.CurrentPlayer(), want[i])
		}
		if err := g.Move(m[0], m[1]); err != nil {
			t.Fatalf("move %d: %v", i, err)
		}
	}
}

func TestNewGameRejectsBadInput(t *testing.T) {
	x := NewPlayer("Xavier", X)
	o := NewPlayer("Olivia", O)

	tests := []struct {
		name      string
		size      int
		winLength int
		players   []*Player
		want      error
	}{
		{"one player", 3, 3, []*Player{x}, ErrTooFewPlayers},
		{"nil player", 3, 3, []*Player{x, nil}, ErrNilPlayer},
		{"duplicate signs", 3, 3, []*Player{x, NewPlayer("Other", X)}, ErrDuplicateSymbol},
		{"empty sign", 3, 3, []*Player{x, NewPlayer("Ghost", Empty)}, ErrEmptySymbol},
		{"zero size", 0, 3, []*Player{x, o}, ErrInvalidSize},
		{"negative size", -5, 3, []*Player{x, o}, ErrInvalidSize},
		{"win length exceeds board", 3, 4, []*Player{x, o}, ErrInvalidWinLength},
		{"zero win length", 3, 0, []*Player{x, o}, ErrInvalidWinLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGame(tt.size, tt.winLength, tt.players...)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
			if g != nil {
				t.Fatal("a rejected game must not be returned")
			}
		})
	}
}

func TestBoardStartsEmpty(t *testing.T) {
	g, _, _ := newTestGame(t)
	b := g.Board()

	for row := 0; row < b.Size(); row++ {
		for col := 0; col < b.Size(); col++ {
			cell, err := b.GetCell(row, col)
			if err != nil {
				t.Fatalf("GetCell(%d,%d): %v", row, col, err)
			}
			if cell != Empty {
				t.Errorf("cell (%d,%d) = %v, want Empty", row, col, cell)
			}
		}
	}

	if b.IsFull() {
		t.Error("a fresh board must not report full")
	}
}

func TestGetCellOutOfBounds(t *testing.T) {
	g, _, _ := newTestGame(t)

	if _, err := g.Board().GetCell(3, 3); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("err = %v, want ErrOutOfBounds", err)
	}
}

func TestBoardString(t *testing.T) {
	g, _, _ := newTestGame(t)
	playAll(t, g, [][2]int{{0, 0}, {1, 1}})

	want := "X - -\n- O -\n- - -\n"
	if got := g.Board().String(); got != want {
		t.Fatalf("String() =\n%q\nwant\n%q", got, want)
	}
}
