package main

import (
	"flag"
	"log"
	"time"

	tb "gopkg.in/telebot.v3"

	"reminder-bot/db"
	"reminder-bot/handlers"
	"reminder-bot/state"
)

func mustToken() string {
	token := flag.String("token", "", "токен телеграмм бота")
	flag.Parse()

	if *token == "" {
		log.Fatal("токена не обнаружено")
	}

	return *token
}

func main() {
	botToken := mustToken()

	pref := tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Авторизован как %s", bot.Me.Username)

	err = db.InitDB()
	if err != nil {
		log.Print("база данных не активирована", err)
	}

	go db.StartReminderWorker(bot)

	menu := &tb.ReplyMarkup{ResizeKeyboard: true}
	btnAdd := menu.Text("Добавить напоминание")
	menu.Reply(menu.Row(btnAdd))

	menuInline := &tb.ReplyMarkup{}

	repeatBtnY := menuInline.Data(
		"да",
		"add_repeat",
	)

	repeatBtnN := menuInline.Data("нет", "not_add_repeat")

	bot.Handle("/start", handlers.StartHandler)
	bot.Handle("/ping", handlers.PingHandler)
	bot.Handle("ping", handlers.PingHandler)
	bot.Handle(&repeatBtnY, func(c tb.Context) error {
		userID := c.Sender().ID

		flow, ok := state.Flows[userID]
		if !ok {
			return c.Send("Ошибка состояния")
		}

		flow.State = state.StateAddRepeatInterval

		// handlers.HandlerAddRepeatInterval(c, flow)

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
	})
	bot.Handle(&repeatBtnN, func(c tb.Context) error {
		userID := c.Sender().ID

		flow, ok := state.Flows[userID]
		if !ok {
			return c.Send("Ошибка состояния")
		}

		err := handlers.HandlerAddNewRemindWithoutRepeat(c, flow)
		if err != nil {
			return c.Send("Ошибка сохранения")
		}

		delete(state.Flows, userID)

		return c.Send("Ок, напоминание без повторения")
	})
	bot.Handle("Добавить напоминание", func(c tb.Context) error {
		userID := c.Sender().ID

		state.Flows[userID] = &state.UserFlow{
			State: state.StateAddText,
		}

		return c.Send("Введи текст напоминания")
	})
	bot.Handle(tb.OnText, func(c tb.Context) error {
		if c.Sender().IsBot {
			log.Print("penis")
			return nil
		}
		userID := c.Sender().ID

		flow, ok := state.Flows[userID]
		if !ok {
			return c.Send("Не понимаю(")
		}

		switch flow.State {

		case state.StateAddText:
			return handlers.HandleAddText(c, flow)

		case state.StateAddTime:
			return handlers.HandleAddTime(c, flow)

		case state.StateAddRepeatInterval:
			return handlers.HandlerAddRepeatInterval(c, flow)
		}

		return nil
	})

	log.Println("Бот запущен...")
	bot.Start()
}
