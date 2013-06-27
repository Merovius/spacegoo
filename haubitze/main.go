package haubitze

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
	"math/rand"
)

type Haubitze struct {
	Victim Planet
}

func chooseVictim(state GameState) Planet {
	NotMine := state.NotMyPlanets()

	if len(NotMine) == 0 {
		n := rand.Intn(len(state.Planets))
		return state.Planets[n]
	}

	n := rand.Intn(len(NotMine))
	return NotMine[n]
}

func (bot *Haubitze) Move(state GameState) Move {
	if state.Round%10 == 0 {
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
	master.Register("haubitze", &Haubitze{})
}
