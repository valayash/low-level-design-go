package tictactoe

// Player is a participant in a game. It carries no game state of its own —
// whose turn it is belongs to Game, not here.
type Player struct {
	Name string
	Sign Symbol
}

// NewPlayer builds a Player. It is a plain function, not a method: a
// constructor cannot require an instance of the thing it constructs.
func NewPlayer(name string, sign Symbol) *Player {
	return &Player{Name: name, Sign: sign}
}

func (p *Player) String() string {
	return p.Name
}
