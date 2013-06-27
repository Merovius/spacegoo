package master

import (
	"github.com/Merovius/spacegoo"
	"log"
)

var (
	bots   = make(map[string]spacegoo.Bot)
)

func Run(name string, server string, user string, pass string) {
	if bot, ok := bots[name]; ok {
		spacegoo.Run(bot, server, user, pass)
		return
	}
	log.Fatal("No such bot")
}

func Register(name string, bot spacegoo.Bot) {
	if bot == nil {
		panic("bot is nil")
	}
	if _, dup := bots[name]; dup {
		panic("name already taken")
	}
	bots[name] = bot
}
