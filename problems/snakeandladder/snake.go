package snakeandladder

type Snake struct {
	start int
	end   int
}

func NewSnake(start, end int) *Snake {
	return &Snake{
		start: start,
		end:   end,
	}
}
