package handlers

import (
	"log"

	tb "gopkg.in/telebot.v3"
)

func sendErr(c tb.Context, err error) error {
	if err != nil {
		log.Println(err)
		return c.Send("Ошибка 😢")
	}
	return nil
}
