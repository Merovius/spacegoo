package main

import (
	"flag"
	"github.com/Merovius/spacegoo"
)

var (
	server = flag.String("server", "spacegoo.gpn.entropia.de:6000", "server address (host:port)")
)

type NopBot struct{}

func (bot *NopBot) Move(state spacegoo.GameState) (spacegoo.Move, error) {
	return state.Nop(), nil
}

func main() {
	flag.Parse()
	username := flag.Arg(0)
	password := flag.Arg(1)

	spacegoo.Run(&NopBot{}, *server, username, password)
}
