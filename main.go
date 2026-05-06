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

	bot.Handle("/start", handlers.StartHandler)
	bot.Handle("/ping", handlers.PingHandler)
	bot.Handle("ping", handlers.PingHandler)
	bot.Handle("Добавить напоминание", func(c tb.Context) error {
		userID := c.Sender().ID

		state.Flows[userID] = &state.UserFlow{
			State: state.StateAddText,
		}

		return c.Send("Введи текст напоминания")
	})
	bot.Handle(tb.OnText, func(c tb.Context) error {
		userID := c.Sender().ID

		flow, ok := state.Flows[userID]
		if !ok {
			return c.Send("Не понимаю")
		}

		switch flow.State {

		case state.StateAddText:
			return handlers.HandleAddText(c, flow)

		case state.StateAddTime:
			return handlers.HandleAddTime(c, flow)
		}

		return nil
	})

	log.Println("Бот запущен...")
	bot.Start()
}
