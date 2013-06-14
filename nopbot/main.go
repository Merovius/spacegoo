package main

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/boilerplate"
)

type NopBot struct{}

func (bot *NopBot) Move(state GameState) (Move, error) {
	return state.Nop(), nil
}

func main() {
	boilerplate.Run(&NopBot{})
}
