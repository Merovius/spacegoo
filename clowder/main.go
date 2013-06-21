package main

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/boilerplate"
)

type Clowder struct {
	init   bool
	Target Planet
}

func (bot *Clowder) Move(state GameState) Move {
	mine := state.MyPlanets()
	if len(mine) == 0 {
		return Nop{}
	}
	bot.Target = state.Planets[bot.Target.Id]

	if !bot.init || bot.Target.Owner == We {
		X, Y := mine.Center()
		notmine := state.NotMyPlanets()
		notmine = notmine.SortByFDist(X, Y)
		if len(notmine) == 0 {
			return Nop{}
		}
		bot.Target = notmine[0]
		bot.init = true
	}

	mine = mine.SortByShips()
	p := mine[len(mine)-1]

	return Send{p, bot.Target, p.Ships}
}

func main() {
	boilerplate.Run(&Clowder{})
}
