package spacegoo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

/* A single planet as reported by the server */
type Planet struct {
	X int `json:"x"`
	Y int `json:"y"`
	Production []int `json:"production"`
	Ships []int `json:"ships"`
	OwnerId int `json:"owner_id"`
	Id int `json:"id"`
}

/* The complete gamestate, we get each round */
type GameState struct {
	Round int `json:"round"`
	MaxRounds int `json:"max_rounds"`
	GameOver bool `json:"game_over"`
	PlayerId int `json:"player_id"`
	// TODO: fleets
	// TODO: players
	Planets []Planet `json:"planets"`
}

type Game struct {
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

	return &Game{ r }, nil
}

func (g *Game) Next() (*GameState, error) {
	line, err := g.r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	for {
		if !strings.HasPrefix(line, "{") {
			if strings.HasPrefix(line, "your current score:") {
				log.Println(line)
				continue
			} else if strings.HasPrefix(line, "game starts.") {
				log.Println(line)
				continue
			} else {
				return nil, fmt.Errorf("unhandled: %s\n", line)
			}
		}
		var state GameState
		// decode the json
		log.Printf("parsing line %s\n", line)
		err = json.Unmarshal([]byte(line), &state)
		if err != nil {
			return nil, err
		}

		log.Printf("state received: %v\n", state)
		return &state, nil
	}
}