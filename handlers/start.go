package handlers

import (
	"log"
	db "reminder-bot/db"

	tb "gopkg.in/telebot.v3"
)

func StartHandler(c tb.Context) error {
	err := db.SaveUser(c.Sender())
	if err != nil {
		log.Fatal("ошибка при сохранении в бд")
	}
	menu := &tb.ReplyMarkup{ResizeKeyboard: true}
	btnAdd := menu.Text("Добавить напоминание")
	menu.Reply(menu.Row(btnAdd))
	return c.Send("Давайте начнём", menu)
}
