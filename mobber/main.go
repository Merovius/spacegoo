package mobber

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
	"math/rand"
)

type Mobber struct {
	init   bool
	Victim Planet
}

func chooseVictim(state GameState) Planet {
	NotMine := state.NotMyPlanets()
	NotMine = NotMine.SortByShips()
	if len(NotMine) == 0 {
		n := rand.Intn(len(state.Planets))
		return state.Planets[n]
	}

	return NotMine[0]
}

func (bot *Mobber) Move(state GameState) Move {
	if !bot.init {
		bot.Victim = chooseVictim(state)
		bot.init = true
	}
	bot.Victim = state.Planets.Lookup(bot.Victim.Id)
	if bot.Victim.Owner == We || state.Round%10 == 0 {
		bot.Victim = chooseVictim(state)
	}

	Mine := state.MyPlanets()
	Mine = Mine.SortByShips()

	if len(Mine) == 0 {
		return Nop{}
	}
	p := Mine[len(Mine)-1]

	return Send{p, bot.Victim, p.Ships}
}

func init() {
	master.Register("mobber", &Mobber{})
}
