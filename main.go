package main

import (
	"flag"
	"log"
	"time"

	tb "gopkg.in/telebot.v3"

	"reminder-bot/handlers"
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

	// регистрация хендлеров
	bot.Handle("/start", handlers.StartHandler)
	bot.Handle("/ping", handlers.PingHandler)
	bot.Handle("ping", handlers.PingHandler)

	log.Println("Бот запущен...")
	bot.Start()
}
