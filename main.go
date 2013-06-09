package spacegoo

import (
	"fmt"
	"math"
	"sort"
)

type Ships [3]int
type fShips [3]float64
type Fleets []Fleet

type Player int

const (
	Neutral Player = iota
	We
	They
)

type Bot interface {
	Move(GameState) (Move, error)
}

/* The complete gamestate, we get each round */
type GameState struct {
	Round     int
	MaxRounds int
	GameOver  bool
	Fleets    Fleets
	Planets   []Planet
	pid       int
}

/* A fleet */
type Fleet struct {
	Id     int
	Owner  Player
	Origin Planet
	Target Planet
	Ships  Ships
	Eta    int
}

/* A single planet as reported by the server */
type Planet struct {
	X          int
	Y          int
	Production Ships
	Ships      Ships
	Owner      Player
	Id         int
}

func battleRound(mine, other fShips) fShips {
	for def_type := 0; def_type < 3; def_type += 1 {
		for att_type := 0; att_type < 3; att_type += 1 {
			c := (def_type - att_type) % 3
			var multiplier float64
			var absolute float64
			switch c {
			case 0:
				multiplier = 0.1
				absolute = 1
			case 1:
				multiplier = 0.25
				absolute = 2
			case 2:
				multiplier = 0.01
				absolute = 1
			}

			var more float64
			if mine[att_type] > 0 {
				more = 1
			} else {
				more = 0
			}
			other[def_type] -= (float64(mine[att_type]) * multiplier) + more*absolute
		}
		if other[def_type] < 0 {
			other[def_type] = 0
		}
	}
	return other
}

func (s Ships) Sum() int {
	return s[0] + s[1] + s[2]
}

func (s Ships) float() (f fShips) {
	for i := 0; i < 2; i++ {
		f[i] = float64(s[i])
	}
	return
}

func (f fShips) Sum() float64 {
	return f[0] + f[1] + f[2]
}

func (s1 Ships) Add(s2 Ships) (s3 Ships) {
	for i := 0; i < 2; i++ {
		s3[i] = s1[i] + s2[i]
	}
	return
}

func (f fShips) Ships() (s Ships) {
	for i := 0; i < 2; i++ {
		s[i] = int(f[i])
	}
	return
}

func (s Ships) String() string {
	return fmt.Sprintf("(%d %d %d)", s[0], s[1], s[2])
}

func Simulate(mine, other Ships) (minenew, othernew Ships) {
	mineS := mine.float()
	otherS := other.float()

	for mineS.Sum() > 0 && otherS.Sum() > 0 {
		new1 := battleRound(otherS, mineS)
		otherS = battleRound(mineS, otherS)
		mineS = new1
	}

	return mineS.Ships(), otherS.Ships()
}

func (s *GameState) MyPlanets() (my []Planet) {
	for _, p := range s.Planets {
		if p.Owner == We {
			my = append(my, p)
		}
	}
	return
}

func (s *GameState) NotMyPlanets() (theirs []Planet) {
	for _, p := range s.Planets {
		if p.Owner != We {
			theirs = append(theirs, p)
		}
	}
	return
}

func (s *GameState) TheirPlanets() (enem []Planet) {
	for _, p := range s.Planets {
		if p.Owner == They {
			enem = append(enem, p)
		}
	}
	return
}

func (s *GameState) NeutralPlanets() (neutr []Planet) {
	for _, p := range s.Planets {
		if p.Owner == Neutral {
			neutr = append(neutr, p)
		}
	}
	return
}

func (p1 *Planet) Dist(x, y float64) float64 {
	dx := float64(p1.X) - x
	dy := float64(p1.Y) - y
	return math.Sqrt(dx*dx + dy*dy)
}

type fleetSorter struct {
	f Fleets
}

func (fs fleetSorter) Len() int {
	return len(fs.f)
}

func (fs fleetSorter) Less(i, j int) bool {
	return fs.f[i].Eta < fs.f[j].Eta
}

func (fs fleetSorter) Swap(i, j int) {
	t := fs.f[i]
	fs.f[j] = fs.f[i]
	fs.f[i] = t
}

func (f Fleets) Sort() {
	sort.Sort(fleetSorter{f})
}
