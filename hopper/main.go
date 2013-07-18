package hopper

import (
	. "github.com/Merovius/spacegoo"
	"github.com/Merovius/spacegoo/master"
	"math/rand"
)

type State func(GameState) (Move, State)

type Hopper struct {
	State   State
	Home    Planet
	NewHome Planet
}

func (bot *Hopper) sendRandomHome(state GameState) Move {
	mine := state.MyPlanets()
	if len(mine) == 1 {
		return Nop{}
	}
	for {
		n := rand.Intn(len(mine))
		if mine[n].Id == bot.Home.Id {
			continue
		}
		return Send{mine[n], bot.Home, mine[n].Ships}
	}
}

func (bot *Hopper) begin(state GameState) (Move, State) {
	mine := state.MyPlanets()
	n := rand.Intn(len(mine))
	bot.Home = mine[n]

	return bot.sendRandomHome(state), bot.awaitingAttack
}

func (bot *Hopper) awaitingAttack(state GameState) (Move, State) {
	incoming := state.Attacking(bot.Home)

	if len(incoming) > 0 {
		return bot.chooseHome(state)
	}

	return bot.sendRandomHome(state), bot.awaitingAttack
}

func (bot *Hopper) chooseHome(state GameState) (Move, State) {
	incoming := state.Attacking(bot.Home)
	first := incoming[0]
	roundsUntilAttack := first.Eta - state.Round

	var canBeat Planets
	for _, p := range state.Planets {
		if p.Id == bot.Home.Id {
			continue
		}
		if p.Owner == We {
			canBeat = append(canBeat, p)
			continue
		}

		my := bot.Home.Ships
		th := p.Ships
		dist := bot.Home.Dist(p.X, p.Y)
		if p.Owner != Neutral {
			th = th.Add(p.Production.Scale(float64(dist + 1 + roundsUntilAttack)))
		}
		newmy, _ := Simulate(my, th)
		if newmy.Sum() > 0 {
			canBeat = append(canBeat, p)
		}
	}

	if len(canBeat) == 0 {
		maxDist := 0
		var target Planet
		for _, p := range state.Planets {
			dist := p.Dist(bot.Home.X, bot.Home.Y)
			if dist > maxDist {
				maxDist = dist
				target = p
			}
		}
		canBeat = append(canBeat, target)
	}

	n := rand.Intn(len(canBeat))
	bot.NewHome = canBeat[n]
	return bot.incoming(state)
}

func (bot *Hopper) incoming(state GameState) (Move, State) {
	mine := state.MyPlanets()

	incoming := state.Attacking(bot.Home)
	incoming.Sort()

	if len(incoming) == 0 {
		return Nop{}, bot.awaitingAttack
	}
	first := incoming[0]

	if first.Eta == state.Round {
		return Send{bot.Home, bot.NewHome, bot.Home.Ships}, bot.migrate
	}

	if len(mine) == 1 {
		return Nop{}, bot.incoming
	}

	mine = mine.SortByShips()
	for i := len(mine) - 1; i >= 0; i-- {
		if mine[i].Id == bot.Home.Id {
			continue
		}
		if bot.Home.Dist(mine[i].X, mine[i].Y)+state.Round < first.Eta {
			return Send{mine[i], bot.Home, mine[i].Ships}, bot.incoming
		}
		if bot.NewHome.Dist(mine[i].X, mine[i].Y)+state.Round >=
			first.Eta+bot.NewHome.Dist(bot.Home.X, bot.Home.Y) {
			return Send{mine[i], bot.NewHome, mine[i].Ships}, bot.incoming
		}
	}
	return Send{mine[0], bot.NewHome, mine[0].Ships}, bot.incoming
}

func (bot *Hopper) migrate(state GameState) (Move, State) {
	if bot.NewHome.Owner == We {
		bot.Home = bot.NewHome
		return bot.awaitingAttack(state)
	}

	mine := state.MyPlanets()
	if len(mine) == 0 {
		return Nop{}, bot.migrate
	}

	var invading Fleet
	for _, f := range state.Incoming(bot.NewHome) {
		if f.Owner == We && f.Origin.Id == bot.Home.Id {
			invading = f
		}
	}

	mine.SortByShips()
	for i := len(mine) - 1; i >= 0; i-- {
		if bot.NewHome.Dist(mine[i].X, mine[i].Y) >= invading.Eta {
			return Send{mine[i], bot.NewHome, mine[i].Ships}, bot.migrate
		}
	}

	return Send{mine[0], bot.NewHome, mine[0].Ships}, bot.migrate
}

func (bot *Hopper) Move(state GameState) Move {
	bot.Home = state.Planets.Lookup(bot.Home.Id)
	bot.NewHome = state.Planets.Lookup(bot.NewHome.Id)

	m, s := bot.State(state)
	bot.State = s
	return m
}

func init() {
	bot := &Hopper{}
	bot.State = bot.begin
	master.Register("hopper", bot)
}
