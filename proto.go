package spacegoo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type Move interface {
	String() string
}

type nop bool
type send struct {
	origin int
	target int
	fleet  Ships
}

type rawGameState struct {
	Round     int         `json:"round"`
	MaxRounds int         `json:"max_rounds"`
	GameOver  bool        `json:"game_over"`
	PlayerId  int         `json:"player_id"`
	Fleets    []rawFleet  `json:"fleets"`
	Planets   []rawPlanet `json:"planets"`
	// TODO: players
}

type rawFleet struct {
	Id    int   `json:"id"`
	Owner int   `json:"owner"`
	Oid   int   `json:"origin"`
	Tid   int   `json:"target"`
	Ships Ships `json:"ships"`
	Eta   int   `json:"eta"`
}

type rawPlanet struct {
	X          int   `json:"x"`
	Y          int   `json:"y"`
	Production Ships `json:"production"`
	Ships      Ships `json:"ships"`
	OwnerId    int   `json:"owner_id"`
	Id         int   `json:"id"`
}

func (rp rawPlanet) Nice(s *GameState) (p Planet) {
	p.X = rp.X
	p.Y = rp.Y
	p.Production = rp.Production
	p.Ships = rp.Ships

	switch rp.OwnerId {
	case 0:
		p.Owner = Neutral
	case s.pid:
		p.Owner = We
	case 3 - s.pid:
		p.Owner = They
	}

	p.Id = rp.Id
	return
}

func (rf rawFleet) Nice(s *GameState) (f Fleet) {
	f.Id = rf.Id

	switch rf.Owner {
	case 0:
		f.Owner = Neutral
	case s.pid:
		f.Owner = We
	case 3 - s.pid:
		f.Owner = They
	}

	f.Origin = s.Planets[rf.Oid]
	f.Target = s.Planets[rf.Tid]
	f.Ships = rf.Ships
	f.Eta = rf.Eta
	return
}

func (rs *rawGameState) Nice() *GameState {
	s := &GameState{}
	s.Round = rs.Round
	s.MaxRounds = rs.MaxRounds
	s.GameOver = rs.GameOver
	s.pid = rs.PlayerId

	for _, rp := range rs.Planets {
		s.Planets = append(s.Planets, rp.Nice(s))
	}

	for _, rf := range rs.Fleets {
		s.Fleets = append(s.Fleets, rf.Nice(s))
	}
	return s
}

func (n nop) String() string {
	return "nop\n"
}

func (s send) String() string {
	return fmt.Sprintf("%d %d %d %d %d", s.origin, s.target, s.fleet[0], s.fleet[1], s.fleet[2])
}

// Send "type1" ships of type 1… from planet "from" to planet "to"
func (s *GameState) Send(from Planet, to Planet, fleet Ships) Move {
	return send{from.Id, to.Id, fleet}
}

// Do nothing
func (s *GameState) Nop() Move {
	return nop(true)
}

func Run(bot Bot, server string, user string, pass string) error {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return err
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	_, err = r.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("logging in with user %s, pass %s\n", user, pass)
	fmt.Fprintf(conn, "login %s %s\n", user, pass)

	for {
		state, err := nextState(conn, r)
		if err != nil {
			return err
		}
		move, err := bot.Move(*state)
		if err != nil {
			return err
		}
		fmt.Fprintf(conn, "%s", move.String())
	}

	return nil
}

func nextState(c net.Conn, r *bufio.Reader) (*GameState, error) {
	state := &rawGameState{}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(line, "{") {
			if strings.HasPrefix(line, "your current score:") {
				log.Printf("%s", line)
				continue
			} else if strings.HasPrefix(line, "game starts.") {
				log.Printf(line)
				continue
			} else if strings.Contains(line, "please disconnect") {
				log.Printf("%s\n", line)
				return nil, fmt.Errorf("disconnected")
			} else {
				log.Printf("unhandled: %s\n", line)
				continue
			}
		}

		// decode the json
		err = json.Unmarshal([]byte(line), state)
		if err != nil {
			return nil, err
		}
		return state.Nice(), nil
	}
}
