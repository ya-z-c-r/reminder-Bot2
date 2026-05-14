package handlers

import (
	"log"
	"reminder-bot/db"
	"reminder-bot/ui"

	tb "gopkg.in/telebot.v3"
)

func StartHandler(c tb.Context) error {
	err := db.SaveUser(c.Sender())
	if err != nil {
		log.Fatal("ошибка при сохранении в бд")
	}
	return c.Send(
		"Давайте начнём",
		ui.MainMenu,
	)
}
