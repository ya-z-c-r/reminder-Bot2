package handlers

import (
	"log"
	"reminder-bot/db"
	"reminder-bot/state"
	"reminder-bot/utils"
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
	//userID := c.Sender().ID

	t, err := utils.ParseHumanTime(c.Text())
	// log.Println("Parsed time:", t, err)
	if err != nil {
		return c.Send("Неверный формат 😢")
	} else if t.Before(time.Now()) {
		return c.Send("указанная дата в прошлом")
	}

	flow.RemindAt = t

	err = SaveNewRimind(c, flow)

	if err != nil {
		c.Send("ошибка сохранения")
		log.Print("ошибка сохранения одноразового напоминания", err)
		return err
	}
	return err
}

func HandlerAddRepeatInterval(c tb.Context, flow *state.UserFlow) error {
	r, err := utils.ParseToCron(c.Text())

	if err != nil {
		log.Println("ошибка при получениее cron", err)
		c.Send("непонял, введите ещё раз")
		for err != nil {
			r, err = utils.ParseToCron(c.Text())
		}
	}

	flow.RepeatInterval = r

	return SaveNewRimind(c, flow)
}


func HandlerAddRepeatText(c tb.Context, flow *state.UserFlow) error {
	flow.Text = c.Text()
	flow.State = state.StateAddRepeatInterval
	return c.Send(`введите период напоминаний. Примеры того как можно вводить:
	каждый день в 15:00
	каждый день в 9:00 и 18:00

	каждый пн в 10:00
	каждый понедельник и четверг в 12:30

	каждые 2 часа
	каждые 15 минут
	каждые 3 дня в 10:00

	каждый месяц 5 числа в 12:00
	каждый месяц 1 и 15 числа в 09:00

	в будни в 10:00
	по выходным в 12:00`)
}

func SaveNewRimind(c tb.Context, flow *state.UserFlow) error {
	userID := c.Sender().ID

	reminder := db.Reminder{
		UserID:         userID,
		Text:           flow.Text,
		RemindAt:       flow.RemindAt,
		RepeatInterval: flow.RepeatInterval,
	}

	err := db.NewRemind(reminder)

	if err != nil {
		delete(state.Flows, userID)
		return c.Send("Ошибка сохранения")
	}

	if reminder.RepeatInterval != "" {
		db.NewReminderChan <- reminder
	}

	delete(state.Flows, userID)

	return c.Send("напоминание создано")
}
