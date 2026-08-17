package snakeandladder

type Ladder struct {
	start int
	end   int
}

func NewLadder(start, end int) *Ladder {
	return &Ladder{
		start: start,
		end:   end,
	}
}
