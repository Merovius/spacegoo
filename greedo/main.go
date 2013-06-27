package greedo

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
)

type Greedo struct{}

func (bot *Greedo) Move(state GameState) Move {
	mine := state.MyPlanets()
	if len(mine) == 0 {
		return Nop{}
	}

	notmine := state.TheirPlanets()

	var max Ships
	var best Planet
	var orig Planet
	for _, p := range notmine {
		if p.Production.Sum() <= max.Sum() {
			continue
		}

		for _, mp := range mine {
			nmy, _ := Simulate(mp.Ships, p.Ships)
			if nmy.Sum() == 0 {
				continue
			}
			max = p.Production
			best = p
			orig = mp
		}
	}

	if max.Sum() > 0 {
		return Send{orig, best, orig.Ships}
	}

	notmine = state.NeutralPlanets()
	for _, p := range notmine {
		if p.Production.Sum() <= max.Sum() {
			continue
		}

		for _, mp := range mine {
			nmy, _ := Simulate(mp.Ships, p.Ships)
			if nmy.Sum() == 0 {
				continue
			}
			max = p.Production
			best = p
			orig = mp
		}
	}

	if max.Sum() > 0 {
		return Send{orig, best, orig.Ships}
	}
	return Nop{}
}

func init() {
	master.Register("greedo", &Greedo{})
}
