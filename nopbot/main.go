package main

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/boilerplate"
)

type NopBot struct{}

func (bot *NopBot) Move(state GameState) Move {
	return Nop{}
}

func main() {
	boilerplate.Run(&NopBot{})
}
