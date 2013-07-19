// Package spacegoo.boilerplate should make writing spacegoo-bots as trivial as
// possible by providing a minimal boilerplate to just start coding.
//
// Now, a minimal Bot (which does… nothing) is just implemented in 16 LOC (of
// which - sadly - still 12 are just boilerplate…):
//
//	package main
//
//	import (
//	    . "github.com/Merovius/spacegoo"
//	    "github.com/Merovius/spacegoo/boilerplate"
//	)
//
//	type NopBot struct {}
//
//	func (bot *NopBot) Move(state GameState) Move {
//	    return Nop{}, nil
//	}
//
//	func main () {
//	    boilerplate.Run(&NopBot{})
//	}
//
// This gives you a binary which expects to arguments [user] and [pass] and has
// a flag -server to give the server to connect to.
package boilerplate

import (
	"flag"
	"github.com/Merovius/spacegoo"
	"log"
)

var (
	server = flag.String("server", "spacegoo.gpn.entropia.de:6000", "server address (host:port)")
	user   string
	pass   string
)

func init() {
	flag.Parse()
	user = flag.Arg(0)
	pass = flag.Arg(1)
	if user == "" || pass == "" {
		log.Fatal("Expected username/password")
	}
}

func Run(bot spacegoo.Bot) {
	spacegoo.Run(bot, *server, user, pass)
}
