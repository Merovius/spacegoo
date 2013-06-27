package master

import (
	"github.com/Merovius/spacegoo"
	"log"
)

var (
	Bots   = make(map[string]spacegoo.Bot)
)

func Run(name string, server string, user string, pass string) {
	if bot, ok := Bots[name]; ok {
		spacegoo.Run(bot, server, user, pass)
		return
	}
	log.Fatal("No such bot")
}

func Register(name string, bot spacegoo.Bot) {
	if bot == nil {
		panic("bot is nil")
	}
	if _, dup := Bots[name]; dup {
		panic("name already taken")
	}
	Bots[name] = bot
}
