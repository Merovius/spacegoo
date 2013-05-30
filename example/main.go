package main

import (
	"github.com/Merovius/spacegoo"
	"log"
	"flag"
)

var (
	server   = flag.String("server", "spacegoo.gpn.entropia.de:6000", "server address (host:port)")
)

func main() {
	flag.Parse()
	username := flag.Arg(0)
	password := flag.Arg(1)

	game, err := spacegoo.NewGame(*server, username, password)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err := game.Next()
		if err != nil {
			log.Fatal(err)
		}
	}
}

