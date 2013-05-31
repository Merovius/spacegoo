package spacegoo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"strings"
)

type Ships [3]int
type fShips [3]float64
type Fleets []*Fleet

/* A fleet */
type Fleet struct {
	Id     int    `json:"id"`
	Owner  int    `json:"owner"`
	Origin int    `json:"origin"`
	Target int    `json:"target"`
	Ships  [3]int `json:"ships"`
	Eta    int    `json:"eta"`
}

/* A single planet as reported by the server */
type Planet struct {
	X          int   `json:"x"`
	Y          int   `json:"y"`
	Production Ships `json:"production"`
	Ships      Ships `json:"ships"`
	OwnerId    int   `json:"owner_id"`
	Id         int   `json:"id"`
}

/* The complete gamestate, we get each round */
type GameState struct {
	Round     int       `json:"round"`
	MaxRounds int       `json:"max_rounds"`
	GameOver  bool      `json:"game_over"`
	PlayerId  int       `json:"player_id"`
	Fleets    Fleets    `json:"fleets"`
	Planets   []*Planet `json:"planets"`
	// TODO: players
}

type Game struct {
	c net.Conn
	r *bufio.Reader
}

// Connects to server with username user and password pass */
func NewGame(server string, user string, pass string) (*Game, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(conn)
	_, err = r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	log.Printf("logging in with user %s, pass %s\n", user, pass)
	fmt.Fprintf(conn, "login %s %s\n", user, pass)

	return &Game{conn, r}, nil
}

// Get the next gamestate
func (g *Game) Next() (*GameState, error) {
	for {
		line, err := g.r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(line, "{") {
			if strings.HasPrefix(line, "your current score:") {
				log.Println(line)
				continue
			} else if strings.HasPrefix(line, "game starts.") {
				log.Println(line)
				continue
			} else if strings.Contains(line, "please disconnect") {
				log.Printf("%s\n", line)
				g.c.Close()
				return nil, fmt.Errorf("disconnected")
			} else {
				log.Printf("unhandled: %s\n", line)
				continue
			}
		}
		var state GameState
		// decode the json
		//		log.Printf("parsing line %s\n", line)
		err = json.Unmarshal([]byte(line), &state)
		if err != nil {
			return nil, err
		}

		//		log.Printf("state received: %v\n", state)
		return &state, nil
	}
}

// Send "type1" ships of type 1… from planet "from" to planet "to"
func (g *Game) Send(from *Planet, to *Planet, fleet Ships) {
	fmt.Fprintf(g.c, "send %d %d %d %d %d\n", from.Id, to.Id, fleet[0], fleet[1], fleet[2])
}

// Do nothing
func (g *Game) Nop() {
	fmt.Fprintf(g.c, "nop\n")
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

func (s *GameState) MyPlanets() (my []*Planet) {
	for _, p := range s.Planets {
		if p.OwnerId == s.PlayerId {
			my = append(my, p)
		}
	}
	return
}

func (s *GameState) TheirPlanets() (theirs []*Planet) {
	for _, p := range s.Planets {
		if p.OwnerId != s.PlayerId {
			theirs = append(theirs, p)
		}
	}
	return
}

func (s *GameState) EnemyPlanets() (enem []*Planet) {
	for _, p := range s.Planets {
		if p.OwnerId != s.PlayerId && p.OwnerId != 0 {
			enem = append(enem, p)
		}
	}
	return
}

func (s *GameState) NeutralPlanets() (neutr []*Planet) {
	for _, p := range s.Planets {
		if p.OwnerId == 0 {
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
