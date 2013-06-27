package nopbot

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
)

type NopBot struct{}

func (bot *NopBot) Move(state GameState) Move {
	return Nop{}
}

func init() {
	master.Register("nopbot", &NopBot{})
}
