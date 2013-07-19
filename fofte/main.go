package fofte

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
)

type FoFtE struct{}

func (bot *FoFtE) Fight(state GameState, p Planet) (Move, bool) {
	incoming := state.Incoming(p)
	if len(incoming) == 0 {
		return Nop{}, true
	}
	incoming = incoming.SortByEta()

	my := p.Ships
	var th Ships
	var eta int
	for _, f := range incoming {
		if f.Owner == We {
			my.Add(f.Ships)
		} else {
			th = f.Ships
			eta = f.Eta - state.Round
			break
		}
	}
	my = my.Add(p.Production.Scale(float64(eta)))

	nmy, _ := Simulate(my, th)
	if nmy.Sum() > 0 {
		return Nop{}, true
	}

	if eta == 0 {
		return Nop{}, false
	}

	mine := state.MyPlanets().SortByDist(p.X, p.Y)

	for i := len(mine) - 1; i > 0; i-- {
		if mine[i].Dist(p.X, p.Y) < eta {
			my = my.Add(mine[i].Ships)
			eta -= 1
		}
	}

	nmy, _ = Simulate(my, th)
	if nmy.Sum() == 0 {
		return Nop{}, false
	}

	for i := len(mine) - 1; i > 0; i-- {
		if mine[i].Dist(p.X, p.Y) < eta {
			return Send{mine[i], p, mine[i].Ships}, true
		}
	}
	return Nop{}, true
}

func (bot *FoFtE) Flight(state GameState, p Planet) (Move, bool) {
	for _, mp := range state.MyPlanets().SortByDist(p.X, p.Y) {
		attacking := state.Attacking(mp)
		if len(attacking) == 0 {
			return Send{p, mp, p.Ships}, true
		}
		attacking = attacking.SortByEta()

		if mp.Dist(p.X, p.Y) < attacking[0].Eta-state.Round-1 {
			return Nop{}, true
		} else if mp.Dist(p.X, p.Y) == attacking[0].Eta-state.Round-1 {
			return Send{p, mp, p.Ships}, true
		}
	}

	return Nop{}, false
}

func (bot *FoFtE) Expand(state GameState, safe Planets) Move {
	if len(safe) == 0 {
		return Nop{}
	}

	safe = safe.SortByShips()
	notmine := state.NotMyPlanets().SortByShips()

	if len(notmine) == 0 {
		return Nop{}
	}
	p := safe[len(safe)-1]

	my := p.Ships
	th := notmine[0].Ships
	dist := notmine[0].Dist(p.X, p.Y)
	th = th.Add(notmine[0].Production.Scale(float64(dist)))

	nmy, _ := Simulate(my, th)
	if nmy.Sum() == 0 {
		return Nop{}
	}

	return Send{p, notmine[0], p.Ships}
}

func (bot *FoFtE) Move(state GameState) Move {
	//nop := Nop{}
	var safe Planets
	for _, mp := range state.MyPlanets() {
		attacking := state.Attacking(mp)
		if len(attacking) == 0 {
			safe = append(safe, mp)
			continue
		}

		m, ok := bot.Fight(state, mp)
		if ok {
			//			if m == nop {
			//				break
			//			}
			return m
		}

		m, ok = bot.Flight(state, mp)
		if ok {
			//			if m == nop {
			//				break
			//			}
			return m
		}
	}
	return bot.Expand(state, safe)
}

func init() {
	master.Register("FoFtE", &FoFtE{})
}
