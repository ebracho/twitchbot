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
	client := twitchbot.New(user, []string{channel}, token)
	client.RegisterHandler(reHey, func(channel, text string, c *twitchbot.Client) {
		c.MessageChannel(channel, "sup")
	})
	if err := client.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
