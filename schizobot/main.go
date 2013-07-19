package schizobot

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"

	"math/rand"
)

type Schizobot struct{}

func isTargetedBy(state GameState, my Planet, his Planet) bool {
	for _, fleet := range state.Fleets {
		if fleet.Target.Id == his.Id && fleet.Origin.Id == my.Id {
			return true
		}
	}
	return false
}

func (bot *Schizobot) Expand(state GameState) (Move, float64) {
	mi := float64(len(state.MyPlanets()))
	th := float64(len(state.TheirPlanets()))
	var rate float64
	if mi+th == 0 {
		rate = 1.0
	} else {
		rate = th / (mi + th)
	}

	for _, tp := range state.NotMyPlanets().SortByShips() {
		th := tp.Ships
		for _, mp := range state.MyPlanets().SortByDist(tp.X, tp.Y) {
			if isTargetedBy(state, mp, tp) {
				continue
			}
			my := mp.Ships
			nmy, _ := Simulate(my, th)
			if nmy.Sum() > 0 {
				return Send{mp, tp, mp.Ships}, rate
			}
		}
	}
	return Nop{}, 0.0
}

func (bot *Schizobot) Produce(state GameState) (Move, float64) {
	mi := float64(state.MyProduction().Sum())
	th := float64(state.TheirProduction().Sum())
	var rate float64
	if mi+th == 0 {
		rate = 1.0
	} else {
		rate = th / (mi + th)
	}

	mine := state.MyPlanets()
	if len(mine) == 0 {
		return Nop{}, 0.0
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
		return Send{orig, best, orig.Ships}, rate
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
		return Send{orig, best, orig.Ships}, rate
	}
	return Nop{}, 0.0
}

func (bot *Schizobot) Defend(state GameState) (Move, float64) {
	al := float64(len(state.Planets))
	th := float64(len(state.TheirPlanets()))
	rate := th / al

	fleets := state.Fleets.SortByEta()
	for _, f := range fleets {
		if f.Target.Owner != We {
			continue
		}

		for _, p := range state.MyPlanets().SortByShips().Reverse() {
			if p.Dist(f.Target.X, f.Target.Y) > f.Eta {
				continue
			}
			if isTargetedBy(state, p, f.Target) {
				continue
			}
			return Send{p, f.Target, p.Ships}, rate
		}
	}
	return Nop{}, 0.0
}

func (bot *Schizobot) Move(state GameState) Move {
	r := rand.Float64()
	m, w := bot.Expand(state)
	if r < w {
		return m
	}
	m, w = bot.Defend(state)
	if r < w {
		return m
	}
	return Nop{}
}

func init() {
	master.Register("schizobot", &Schizobot{})
}
