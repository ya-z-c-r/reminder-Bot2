package handlers

import tb "gopkg.in/telebot.v3"

func PingHandler(c tb.Context) error {
	return c.Send("pong")
}
