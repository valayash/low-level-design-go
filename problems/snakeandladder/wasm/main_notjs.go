//go:build !js || !wasm

// This stub exists only so that `go build ./...` and `go vet ./...` succeed on
// a normal host toolchain. Without it the directory would contain no buildable
// Go files and the toolchain would fail with "build constraints exclude all Go
// files". The real entry point is main_js.go.
package main

import "fmt"

func main() {
	fmt.Println("This command is the browser (WebAssembly) build of Snake and Ladder.")
	fmt.Println("Build it with:")
	fmt.Println()
	fmt.Println("  GOOS=js GOARCH=wasm go build -o docs/snakeandladder.wasm ./problems/snakeandladder/wasm")
	fmt.Println()
	fmt.Println("To play on the command line instead, run: go run ./problems/snakeandladder/cmd")
}
