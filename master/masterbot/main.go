package main

import (
	"flag"
	"github.com/Merovius/spacegoo/master"
	"log"

	// List of bots
	_ "github.com/Merovius/spacegoo/clowder"
	_ "github.com/Merovius/spacegoo/expandbot"
	_ "github.com/Merovius/spacegoo/fofte"
	_ "github.com/Merovius/spacegoo/friedzclone"
	_ "github.com/Merovius/spacegoo/greedo"
	_ "github.com/Merovius/spacegoo/haubitze"
	_ "github.com/Merovius/spacegoo/hopper"
	_ "github.com/Merovius/spacegoo/mnemo"
	_ "github.com/Merovius/spacegoo/mobber"
	_ "github.com/Merovius/spacegoo/nopbot"
)

var (
	server = flag.String("server", "spacegoo.gpn.entropia.de:6000", "server address (host:port)")
	user   string
	pass   string
)

func main() {
	flag.Parse()
	name := flag.Arg(0)
	pass := flag.Arg(1)
	if name == "" || pass == "" {
		log.Fatal("Expected botname password")
	}
	log.Println("user:", user, "pass:", pass)
	master.Run(name, *server, name, pass)
}
