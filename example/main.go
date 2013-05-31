package main

import (
	"github.com/Merovius/spacegoo"
	"log"
	"flag"
	"math"
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
		state, err := game.Next()
		if err != nil {
			log.Fatal(err)
		}
		for _, planet := range state.Planets {
			if planet.OwnerId == state.PlayerId {
				log.Printf("My planet: %v\n", planet)
				var nearest *spacegoo.Planet
				var nearestDist = float64(99999999)
				for _, other := range state.Planets {
					if other.Id == planet.Id {
						continue
					}
					if other.Ships[0] + other.Ships[1]  + other.Ships[2] < planet.Ships[0]+ planet.Ships[1]+planet.Ships[2] {
						dist := math.Sqrt(math.Pow(math.Abs(float64(other.X)-float64(planet.X)), 2) + math.Pow(math.Abs(float64(other.Y)-float64(planet.Y)), 2))
						if dist < nearestDist {
							nearestDist = dist
							nearest = other
						}
					}
				}
				if nearest != nil {
					log.Printf("nearest: %d\n", nearest.Id)
					log.Printf("send %d %d %d %d %:\n", planet.Id, nearest, 9999, 9999, 9999)
					game.Send(planet, nearest, spacegoo.Ships { 9999, 9999, 9999 })
				} else {
					game.Nop()
				}

				break
			}
		}
	}
}

