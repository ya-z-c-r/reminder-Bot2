package handlers

import (
	tb "gopkg.in/telebot.v3"
)

func PingHandler(c tb.Context) error {
	// log.Print(c)
	return c.Send("pong")
}
