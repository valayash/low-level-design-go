package snakeandladder

import "fmt"

type SnakeAndLadder struct {
	Board            *Board
	Players          []*Player
	CurrentPlayerIdx int
	Dice             *Dice
	Status           bool
}

func NewGame(names []string) *SnakeAndLadder {

	g := &SnakeAndLadder{
		Board:            NewBoard(),
		Dice:             NewDice(),
		Players:          []*Player{},
		CurrentPlayerIdx: 0,
		Status:           true,
	}

	for _, name := range names {
		g.Players = append(g.Players, NewPlayer(name))
	}

	return g

}

func (g *SnakeAndLadder) Play() {

	fmt.Printf("---- Snake and Ladder: board of %d cells ----\n", g.Board.Size)
	for _, p := range g.Players {
		fmt.Printf("  %s starts at %d\n", p.Name, p.Position)
	}
	fmt.Println()

	turn := 0

	for g.Status {

		currPlayer := g.Players[g.CurrentPlayerIdx]
		dice := g.Dice.Roll()

		turn++
		fmt.Printf("Turn %d: %s rolled %d\n", turn, currPlayer.Name, dice)

		newPosition := g.Board.GetPosition(currPlayer.Position, dice)
		currPlayer.Position = newPosition

		fmt.Printf("  %s is now on %d\n", currPlayer.Name, newPosition)

		if g.Board.IsWin(newPosition) {
			fmt.Printf("\n---- Game Won by %s in %d turns ----\n", currPlayer.Name, turn)
			g.Status = false
		}

		g.CurrentPlayerIdx = (g.CurrentPlayerIdx + 1) % len(g.Players)

	}
}
