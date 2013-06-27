package friedzclone

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
)

type FriedzClone struct{}

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

func init() {
	master.Register("friedzclone", &FriedzClone{})
}
