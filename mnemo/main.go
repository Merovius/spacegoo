package mnemo

import (
	"fmt"
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
	"sort"
)

type Mnemo struct {
	queue Queue
}

type AttackSet struct {
	planets  []Planet
	combined Ships
}

func (set AttackSet) String() string {
	ret := fmt.Sprintf("(%d", set.planets[0].Id)
	for _, p := range set.planets[1:] {
		ret += fmt.Sprintf(" %d", p.Id)
	}
	return ret + ")"
}

func GetAllAttackSets(state GameState) (sets []AttackSet) {
	pl := state.MyPlanets()

	for i := 1; i < (1 << uint(len(pl))); i++ {
		set := AttackSet{}

		nbits := 0
		for j := uint(0); j < uint(len(pl)); j++ {
			if i&(1<<j) != 0 {
				nbits++
			}
		}
		if nbits > 5 {
			continue
		}

		for j := uint(0); j < uint(len(pl)); j++ {
			if i&(1<<j) != 0 {
				set.planets = append(set.planets, pl[j])
				set.combined = set.combined.Add(pl[j].Ships)
			}
		}
		sets = append(sets, set)
	}
	return
}

func (bot *Mnemo) isTargetedByUs(state GameState, his Planet) bool {
	for _, fleet := range state.Fleets {
		if fleet.Target.Id == his.Id && fleet.Owner == We {
			return true
		}
	}

	for _, m := range bot.queue {
		m, ok := m.(Send)
		if !ok {
			continue
		}
		if m.Target.Id == his.Id {
			return true
		}
	}

	return false
}

func (bot *Mnemo) setIsBlocked(set AttackSet) bool {
	for _, m := range bot.queue {
		m, ok := m.(Send)
		if !ok {
			continue
		}
		for _, p := range set.planets {
			if p.Id == m.Origin.Id {
				return true
			}
		}
	}
	return false
}

func (bot *Mnemo) Queue(set AttackSet, target Planet) (m Move) {
	max := 0
	farthest := Planet{}
	for _, p := range set.planets {
		if p.Dist(target.X, target.Y) > max {
			max = p.Dist(target.X, target.Y)
			farthest = p
		}
	}

	for _, p := range set.planets {
		if p.Id == farthest.Id {
			m = Send{p, target, p.Ships}
			continue
		}
		bot.queue.Insert(Send{p, target, p.Ships}, max-target.Dist(p.X, p.Y))
	}

	return
}

func (bot *Mnemo) Move(state GameState) Move {
	m := bot.queue.Shift()
	if m != nil {
		return m
	}

	notmine := state.NotMyPlanets()
	if len(notmine) == 0 {
		return Nop{}
	}

	attackSets := attackSetSlice(GetAllAttackSets(state))
	sort.Sort(attackSets)

	for _, tp := range notmine {
		if bot.isTargetedByUs(state, tp) {
			continue
		}
		for _, set := range attackSets {
			if bot.setIsBlocked(set) {
				continue
			}
			th := tp.Ships
			maxrounds := 0
			for _, p := range set.planets {
				if p.Dist(tp.X, tp.Y) > maxrounds {
					maxrounds = p.Dist(tp.X, tp.Y)
				}
			}

			if tp.Owner != Neutral {
				pr := tp.Production.Scale(float64(maxrounds + 2))
				th = th.Add(pr)
			}

			my := set.combined
			nmy, _ := Simulate(my, th)
			if nmy.Sum() > 0 {
				return bot.Queue(set, tp)
			}
		}
	}
	return Nop{}
}

type attackSetSlice []AttackSet

func (s attackSetSlice) Len() int {
	return len(s)
}

func (s attackSetSlice) Less(i, j int) bool {
	if len(s[i].planets) < len(s[j].planets) {
		return true
	}
	if s[i].combined.Sum() < s[j].combined.Sum() {
		return true
	}
	return false
}

func (s attackSetSlice) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func init() {
	master.Register("mnemo", &Mnemo{})
}
