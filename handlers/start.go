package handlers

import (
	"log"
	db "reminder-bot/db"

	tb "gopkg.in/telebot.v3"
)

func StartHandler(c tb.Context) error {
	db.InitDB()
	err := db.SaveUser(c.Sender())
	if err != nil {
		log.Fatal("ошибка при сохранении в бд")
	}
	return c.Send("Давайте начнём")
}
