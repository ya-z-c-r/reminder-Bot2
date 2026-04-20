package handlers

import tb "gopkg.in/telebot.v3"

func StartHandler(c tb.Context) error {
	return c.Send("Давайте начнём")
}
