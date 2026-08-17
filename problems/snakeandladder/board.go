package snakeandladder

import "fmt"

type Board struct {
	Size    int
	Snakes  []*Snake
	Ladders []*Ladder
}

func NewBoard() *Board {
	board := &Board{
		Size:    100,
		Snakes:  []*Snake{},
		Ladders: []*Ladder{},
	}

	board.InitializeBoard()

	return board
}

func (b *Board) InitializeBoard() {

	b.Snakes = append(b.Snakes, NewSnake(16, 6), NewSnake(48, 26), NewSnake(64, 60), NewSnake(93, 73))
	b.Ladders = append(b.Ladders, NewLadder(1, 38), NewLadder(4, 14), NewLadder(9, 31), NewLadder(21, 42),
		NewLadder(28, 84), NewLadder(51, 67), NewLadder(80, 99))

}

func (b *Board) GetPosition(idx int, steps int) int {

	if idx+steps > b.Size {
		fmt.Printf("  needs exactly %d — roll too big, stays on %d\n", b.Size-idx, idx)
		return idx
	}

	idx += steps
	for _, snake := range b.Snakes {
		if snake.start == idx {
			fmt.Printf("  bitten by the snake at %d -> down to %d\n", snake.start, snake.end)
			idx = snake.end
		}
	}

	for _, ladder := range b.Ladders {
		if ladder.start == idx {
			fmt.Printf("  ladder at %d -> up to %d\n", ladder.start, ladder.end)
			idx = ladder.end
		}
	}

	return idx
}

func (b *Board) IsWin(idx int) bool {
	return idx == b.Size
}
