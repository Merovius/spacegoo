package expandbot

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
)

type ExpandBot struct{}
type FriedzClone struct{}

var friedzclone FriedzClone

func init() {
	master.Register("expandbot", &ExpandBot{})
}

func (bot *FriedzClone) Move(state GameState) Move {
	mine := state.MyPlanets()
	if len(mine) == 0 {
		return Nop{}
	}

	mine = mine.SortByShips()
	me := mine[len(mine)-1]

	theirs := state.NotMyPlanets()
	if len(theirs) == 0 {
		return Nop{}
	}

	theirs = theirs.SortByDist(me.X, me.Y)
	they := theirs[0]

	return Send{me, they, me.Ships}
}

func isTargetedBy(state GameState, my Planet, his Planet) bool {
	for _, fleet := range state.Fleets {
		if fleet.Target.Id == his.Id && fleet.Origin.Id == my.Id {
			return true
		}
	}
	return false
}

func (bot *ExpandBot) EndMode(state GameState) Move {
	mine := state.MyPlanets()
	if len(mine) == 0 {
		return Nop{}
	}
	for _, f := range state.Fleets {
		if f.Owner != They {
			continue
		}

		mine = mine.SortByDist(f.Target.X, f.Target.Y)
		for _, mp := range mine {
			if isTargetedBy(state, mp, f.Target) {
				continue
			}
			return Send{mp, f.Target, mp.Ships}
		}
	}
	return Nop{}
}

func (bot *ExpandBot) Move(state GameState) Move {
	if len(state.TheirPlanets()) == 0 {
		return bot.EndMode(state)
	}
	if state.PlayerName(They) == "intercept" {
		return friedzclone.Move(state)
	}

	for _, mp := range state.MyPlanets() {
		my := mp.Ships
		for _, tp := range state.NotMyPlanets() {
			if isTargetedBy(state, mp, tp) {
				continue
			}
			th := tp.Ships
			nmy, _ := Simulate(my, th)
			if nmy.Sum() > 0 {
				return Send{mp, tp, mp.Ships}
			}
		}
	}
	return Nop{}
}
