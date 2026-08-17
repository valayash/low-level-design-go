// Command main runs a game of Snake and Ladder and prints it to stdout.
//
// Usage:
//
//	go run ./problems/snakeandladder/cmd
//	go run ./problems/snakeandladder/cmd Alice Bob Carol
package main

import (
	"os"

	snl "github.com/valayash/low-level-design-go/problems/snakeandladder"
)

func main() {
	names := os.Args[1:]
	if len(names) == 0 {
		names = []string{"Alice", "Bob"}
	}

	game := snl.NewGame(names)
	game.Play()
}
