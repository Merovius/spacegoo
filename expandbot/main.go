package main

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/boilerplate"
)

type ExpandBot struct{}

func isTargetedBy(state GameState, my Planet, his Planet) bool {
	for _, fleet := range state.Fleets {
		if fleet.Target.Id == his.Id && fleet.Origin.Id == my.Id {
			return true
		}
	}
	return false
}

func (bot *ExpandBot) Move(state GameState) (Move, error) {
	for _, mp := range state.MyPlanets() {
		my := mp.Ships
		for _, tp := range state.NotMyPlanets() {
			if isTargetedBy(state, mp, tp) {
				continue
			}
			th := tp.Ships
			nmy, _ := Simulate(my, th)
			if nmy.Sum() > 0 {
				return state.Send(mp, tp, mp.Ships), nil
			}
		}
	}
	return state.Nop(), nil
}

func main() {
	boilerplate.Run(&ExpandBot{})
}
