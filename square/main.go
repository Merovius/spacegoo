package square

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
	"math/rand"
)

type Square struct {
	init   bool
	Victim Planet
}

func chooseVictim(state GameState) Planet {
	Theirs := state.NotMyPlanets()
	if len(Theirs) > 0 {
		Theirs = Theirs.SortByShips()
		return Theirs[0]
	}

	n := rand.Intn(len(state.Planets))
	return state.Planets[n]
}

func (bot *Square) Move(state GameState) Move {
	if !bot.init {
		bot.Victim = chooseVictim(state)
		bot.init = true
	}
	bot.Victim = state.Planets.Lookup(bot.Victim.Id)

	if bot.Victim.Owner == We {
		bot.Victim = chooseVictim(state)
	}

	Mine := state.MyPlanets()
	if len(Mine) == 0 {
		return Nop{}
	}

	n := rand.Intn(len(Mine))
	p := Mine[n]

	return Send{p, bot.Victim, p.Production}
}

func init() {
	master.Register("square", &Square{})
}
