package main

import (
	"log"
	"os"
	"regexp"

	"github.com/ebracho/twitchbot"
)

var (
	reHey = regexp.MustCompile("hey")
)

func main() {
	user := os.Getenv("TWITCH_USER")
	channel := os.Getenv("TWITCH_CHANNEL")
	token := os.Getenv("TWITCH_TOKEN")
	bot := twitchbot.New(user, []string{channel}, token)
	bot.RegisterHandler(reHey, func(m *twitchbot.Message) string {
		return "sup"
	})
	if err := bot.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
