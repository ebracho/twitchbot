package twitchbot

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/go-irc/irc"
)

type Message struct {
	Target string
	Text   string
}

type HandlerFunc func(channel, text string, client *Client)

type Client struct {
	User       string
	Channels   []string
	OAuthToken string

	handlers map[*regexp.Regexp]HandlerFunc

	client    *irc.Client
	clientMux sync.Mutex
}

func New(user string, channels []string, token string) *Client {
	return &Client{
		User:       user,
		OAuthToken: token,
		Channels:   channels,
		handlers:   make(map[*regexp.Regexp]HandlerFunc),
	}
}

func (t *Client) handle(c *irc.Client, m *irc.Message) {
	fmt.Println(m)
	switch m.Command {
	case RPL_WELCOME:
		c.WriteMessage(&irc.Message{
			Command: "JOIN",
			Params:  []string{strings.Join(t.Channels, ",")},
		})
	case "PING":
		c.WriteMessage(&irc.Message{
			Command: "PONG",
			Params:  []string{t.User},
		})
	case "PRIVMSG":
		if len(m.Params) == 0 {
			break
		}
		channel := m.Params[0]
		text := ""
		if len(m.Params) > 1 {
			text = strings.Join(m.Params[1:], " ")
		}
		for expr, handler := range t.handlers {
			if expr.Match([]byte(text)) {
				handler(channel, text, t)
				break
			}
		}
	}
}

func (t *Client) MessageChannel(channel, text string) {
	t.clientMux.Lock()
	defer t.clientMux.Unlock()
	t.client.WriteMessage(&irc.Message{
		Command: "PRIVMSG",
		Params:  []string{channel, text},
	})
}

func (t *Client) RegisterHandler(expr *regexp.Regexp, h HandlerFunc) {
	t.handlers[expr] = h
}

func (t *Client) Run() error {
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
	t.client = client
	return client.Run()
}
