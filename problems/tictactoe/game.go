package tictactoe

// Status is how a game ends, or that it hasn't.
type Status int

const (
	InProgress Status = iota
	Won
	Draw
)

func (s Status) String() string {
	switch s {
	case Won:
		return "won"
	case Draw:
		return "draw"
	default:
		return "in progress"
	}
}

// Game owns the rules: who plays next, when the game ends, and who won. The
// players are held in a slice rather than named fields so that turn rotation
// is one modulo and a third player costs nothing.
type Game struct {
	board     *Board
	players   []*Player
	turnIndex int
	status    Status
	winner    *Player
}

// NewGame refuses to build a game that could never work correctly: too few
// players, a nil player, a player marked Empty, or two players sharing a sign
// (which would make a winning line unattributable).
func NewGame(size, winLength int, players ...*Player) (*Game, error) {
	if len(players) < 2 {
		return nil, ErrTooFewPlayers
	}

	seen := make(map[Symbol]bool, len(players))
	for _, p := range players {
		if p == nil {
			return nil, ErrNilPlayer
		}
		if p.Sign == Empty {
			return nil, ErrEmptySymbol
		}
		if seen[p.Sign] {
			return nil, ErrDuplicateSymbol
		}
		seen[p.Sign] = true
	}

	board, err := NewBoard(size, winLength)
	if err != nil {
		return nil, err
	}

	return &Game{
		board:   board,
		players: append([]*Player(nil), players...), // copy, so the caller cannot reorder ours
		status:  InProgress,
	}, nil
}

// Board exposes the grid for rendering. Marks placed through it would bypass
// turn order, so callers should treat it as read-only and go through Move.
func (g *Game) Board() *Board { return g.board }

// Status reports whether the game is still running, won, or drawn.
func (g *Game) Status() Status { return g.status }

// CurrentPlayer is whoever moves next.
func (g *Game) CurrentPlayer() *Player { return g.players[g.turnIndex] }

// Winner returns the winning player, and false if nobody has won.
func (g *Game) Winner() (*Player, bool) { return g.winner, g.winner != nil }

// Move plays the current player's sign at (row,col).
//
// The order below is the design: the game-over check comes first so a decided
// game rejects further moves; the board validates position and occupancy; the
// win check runs before the draw check so a final mark that completes a line
// is a win, not a draw; and the turn advances only after a valid move that did
// not end the game, so a rejected move leaves the same player to retry.
func (g *Game) Move(row, col int) error {
	if g.status != InProgress {
		return ErrGameOver
	}

	player := g.CurrentPlayer()
	if err := g.board.Mark(row, col, player.Sign); err != nil {
		return err
	}

	if g.board.IsWinningMove(row, col) {
		g.status = Won
		g.winner = player
		return nil
	}

	if g.board.IsFull() {
		g.status = Draw
		return nil
	}

	g.turnIndex = (g.turnIndex + 1) % len(g.players)
	return nil
}
