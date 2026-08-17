//go:build js && wasm

// Command wasm exposes the Snake and Ladder game to the browser.
//
// It compiles the unmodified snakeandladder package to WebAssembly and
// registers a global JS function playGame(namesArray). The game still writes
// its progress to stdout with fmt.Printf; the page captures that by overriding
// the fs.writeSync shim that wasm_exec.js installs.
//
// Build:
//
//	GOOS=js GOARCH=wasm go build -o docs/snakeandladder.wasm ./problems/snakeandladder/wasm
package main

import (
	"syscall/js"

	snl "github.com/valayash/low-level-design-go/problems/snakeandladder"
)

// playGame runs one complete game. args[0], when present, is a JS array of
// player names. An empty or missing list falls back to two default players —
// NewGame divides by len(Players), so the slice must never be empty.
func playGame(this js.Value, args []js.Value) any {
	var names []string

	if len(args) > 0 && args[0].Type() == js.TypeObject {
		arr := args[0]
		for i := 0; i < arr.Length(); i++ {
			names = append(names, arr.Index(i).String())
		}
	}

	if len(names) == 0 {
		names = []string{"Alice", "Bob"}
	}

	snl.NewGame(names).Play()

	return nil
}

func main() {
	js.Global().Set("playGame", js.FuncOf(playGame))

	// Keep the Go runtime alive so the exported function stays callable.
	select {}
}
