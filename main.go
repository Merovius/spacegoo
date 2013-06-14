package spacegoo

import (
	"fmt"
	"math"
	"sort"
)

// Ships designates a tuple of Ships e.g. for fleets, production…
type Ships [3]int

type fShips [3]float64
type Fleets []Fleet

type Player int

// Queue of Moves
type Queue []Move

const (
	Neutral Player = iota
	We
	They
)

// Bot is the primary interface to the API.
//
// You specify in every Round, given a GameState, what you want to do.
// If you return an error, you loose - so better don't do that.
//
// The canonical Moves are GameState.Nop and GameState.Send
type Bot interface {
	Move(GameState) (Move, error)
}

// GameState is the complete State, given each round, pretty much as specified
// in the protocol
type GameState struct {
	Round     int      // The current Gameround
	MaxRounds int      // The maximum number of rounds
	GameOver  bool     // If the game is over
	Fleets    Fleets   // All fleets on their way currently
	Planets   []Planet // All Planets of all Players
	pid       int
	we        string
	they      string
}

// Fleet is a fleet currently on it's way
type Fleet struct {
	Id     int    // The Id of this fleet. Use this for comparisons
	Owner  Player // The Owner
	Origin Planet // Where this fleet comes from
	Target Planet // Where this fleet is going
	Ships  Ships  // How many ships?
	Eta    int    // How long until it arrives
}

// Planet is a planet
type Planet struct {
	Id         int    // The Id of this planet. Use this for comparisons
	X          int    // The X-coordinate
	Y          int    // The Y-coordinate
	Production Ships  // How many ships this planet produces
	Ships      Ships  // How many ships are stationed here
	Owner      Player // Who owns this planet
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

// Sum returns the absolute size of this fleet
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

// Add another fleet to this one
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

// Split up a fraction of this Fleet
func (s1 Ships) Split(fraction float64) (s2 Ships) {
	for i, s := range s1 {
		s2[i] = int(float64(s) * fraction)
	}
	return
}

// String formats the fleet easily readable
func (s Ships) String() string {
	return fmt.Sprintf("(%d %d %d)", s[0], s[1], s[2])
}

// Simulate a battle between two fleets
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

// MyPlanets gives you all Planets under your control
func (s *GameState) MyPlanets() (my []Planet) {
	for _, p := range s.Planets {
		if p.Owner == We {
			my = append(my, p)
		}
	}
	return
}

// NotMyPlanets gives you all neutral and enemy Planets
func (s *GameState) NotMyPlanets() (theirs []Planet) {
	for _, p := range s.Planets {
		if p.Owner != We {
			theirs = append(theirs, p)
		}
	}
	return
}

// TheirPlanets gives you all enemy Planets
func (s *GameState) TheirPlanets() (enem []Planet) {
	for _, p := range s.Planets {
		if p.Owner == They {
			enem = append(enem, p)
		}
	}
	return
}

// NeutralPlanets gives you all neutral Planets
func (s *GameState) NeutralPlanets() (neutr []Planet) {
	for _, p := range s.Planets {
		if p.Owner == Neutral {
			neutr = append(neutr, p)
		}
	}
	return
}

// Dist calculates the distance of the point (x,y) from this Planet
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

// Sort the Fleets ascending by ETA
func (f Fleets) Sort() {
	sort.Sort(fleetSorter{f})
}

// SimulateIncomingFleets takes as a parameter a number of rounds and an planet
// and calculates every fleet, that arrives in less then rounds. It returns the
// Ships remaining on the planet
func (state GameState) SimulateIncomingFleets(p Planet, rounds int) (s Ships) {
	// BUG(Merovius) This almost certainly is not accurate if the planet is
	// neutral or attacked at the same time
	s = p.Ships
	incoming := state.Incoming(p)
	incoming.Sort()
	for i := 1; i < rounds; i++ {
		var attack Ships
		for _, f := range incoming {
			f.Eta -= 1
			if f.Eta == 0 {
				switch f.Owner {
				case p.Owner:
					s.Add(f.Ships)
				default:
					attack.Add(f.Ships)
				}
			}
		}
		if attack.Sum() > 0 {
			defnew, attnew := Simulate(s, attack)
			if defnew.Sum() == 0 {
				s = attnew
				p.Owner = 3 - p.Owner
			} else {
				s = defnew
			}
		}
	}
	return
}

// Incoming fleets to a planet. If no fleets are on their way, returns nil
func (state GameState) Incoming(p Planet) (f Fleets) {
	for _, fleet := range state.Fleets {
		if fleet.Target.Id != p.Id {
			continue
		}
		f = append(f, fleet)
	}
	return f
}

// Return the name of the Player
func (state GameState) PlayerName(p Player) string {
	switch p {
	case Neutral:
		return "Neutral"
	case We:
		return state.we
	case They:
		return state.they
	}
	// never reached
	return ""
}

// Insert puts a move at a specific position, if it is unoccupied.
// If the slot is occupied, it shift the other moves up one position to free it.
// If this is not what you want, test the slot for nil before calling Insert.
func (q *Queue) Insert(m Move, pos int) {
	// If pos exceeds the capacity of the queue, we can reallocate and
	// simpley add the move at the new position
	if cap(*q) < pos {
		nq := make([]Move, 2*pos)
		copy(nq, *q)
		nq[pos] = m
		*q = nq
		return
	}

	// If pos is not occupied yet, we insert the move there
	if (*q)[pos] == nil {
		(*q)[pos] = m
		return
	}

	// pos is occupied. Move the continuous bit up
	for i, v := range (*q)[pos:] {
		if v != nil {
			continue
		}
		// Found a free space, move everything up
		copy((*q)[pos+1:i+1], (*q)[pos:i])
		// Insert the move
		(*q)[pos] = m
		return
	}

	// End of the slice. Can we extend it?
	if len(*q) < cap(*q) {
		copy((*q)[pos+1:len(*q)+1], (*q)[pos:])
		(*q)[pos] = m
		return
	}

	// We have to reallocate
	nq := make([]Move, 2*cap(*q))
	copy(nq, (*q)[:pos])
	nq[pos] = m
	copy(nq[pos+1:], (*q)[pos+1:])

	return
}

// Shift dequeues the first Move and returns it
func (q *Queue) Shift() (m Move) {
	m = (*q)[0]
	(*q) = (*q)[1:]
	return
}
