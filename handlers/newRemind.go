package handlers

import (
	"reminder-bot/db"
	"reminder-bot/state"
	"time"

	tb "gopkg.in/telebot.v3"
)

func HandleAddText(c tb.Context, flow *state.UserFlow) error {
	flow.Text = c.Text()            // сохраняем текст
	flow.State = state.StateAddTime // меняем состояние

	return c.Send(`время напоминания можно вводить: 
15:04
сегодня 18:30
завтра 15:00
завтра в 15:00
послезавтра 12:00

через 10 минут
через 2 часа
через 3 дня
или в формате 2006-01-02 15:04 и 02.01.2006 15:04`)
}

func HandleAddTime(c tb.Context, flow *state.UserFlow) error {
	userID := c.Sender().ID

	t, err := db.ParseHumanTime(c.Text())
	// log.Println("Parsed time:", t, err)
	if err != nil {
		return c.Send("Неверный формат 😢")
	} else if t.Before(time.Now()) {
		return c.Send("указанная дата в прошлом")
	}

	err = db.NewRemind(db.Reminder{
		UserID:   userID,
		Text:     flow.Text,
		RemindAt: t,
	})

	if err != nil {
		delete(state.Flows, userID)
		return c.Send("Ошибка сохранения")
	}

	delete(state.Flows, userID)

	return c.Send("Напоминание создано")
}
