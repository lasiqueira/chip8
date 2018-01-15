package main

import (
	"os"

	"github.com/lasiqueira/chip8/cpu"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Need to inform the game path")
	}
	game := args[0]
	chip8 := cpu.CPU{}

	chip8.Initialize()
	chip8.LoadGame(game)

	return
}
