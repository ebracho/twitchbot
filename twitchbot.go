package twitchbot

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-irc/irc"
)

type Message struct {
	Target string
	Text   string
}

type HandlerFunc func(m *Message) string

type TwitchBot struct {
	User       string
	Channels   []string
	OAuthToken string

	handlers map[*regexp.Regexp]HandlerFunc
}

func New(user string, channels []string, token string) *TwitchBot {
	return &TwitchBot{
		User:       user,
		OAuthToken: token,
		Channels:   channels,
		handlers:   make(map[*regexp.Regexp]HandlerFunc),
	}
}

func (t *TwitchBot) handle(c *irc.Client, m *irc.Message) {
	fmt.Println(m)
	switch m.Command {
	case RPL_WELCOME:
		c.WriteMessage(&irc.Message{
			Command: "JOIN",
			Params:  t.Channels,
		})
	case "PING":
		c.WriteMessage(&irc.Message{
			Command: "PONG",
			Params:  []string{t.User},
		})
	case "PRIVMSG":
		if len(m.Params) < 2 {
			break
		}
		target := m.Params[0]
		text := ""
		if len(m.Params) > 1 {
			text = strings.Join(m.Params[1:], " ")
		}
		for expr, handler := range t.handlers {
			if expr.Match([]byte(text)) {
				response := handler(&Message{
					Target: target,
					Text:   text,
				})
				c.Write(fmt.Sprintf("PRIVMSG %s %s", target, response))
				break
			}
		}
	}
}

func (t *TwitchBot) RegisterHandler(expr *regexp.Regexp, h HandlerFunc) {
	t.handlers[expr] = h
}

func (t *TwitchBot) Run() error {
	conn, err := tls.Dial("tcp", "irc.chat.twitch.tv:6697", &tls.Config{})
	if err != nil {
		return err
	}

	cfg := irc.ClientConfig{
		Nick: t.User,
		Name: t.User,
		User: t.User,
		Pass: t.OAuthToken,
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
			t.handle(c, m)
		}),
	}

	client := irc.NewClient(conn, cfg)
	return client.Run()
}
