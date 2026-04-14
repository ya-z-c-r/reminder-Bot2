package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Замените "YOUR_BOT_TOKEN" на реальный токен от BotFather
	botToken := "8775462367:AAHtFEify4Z_9lnNWQ6Ot_uX73OfwuZwW4s"

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Включаем режим отладки (опционально)
	bot.Debug = true

	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обработка сигнала для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				switch update.Message.Command() {
				case "ping":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
					}
				}
			}
		}
	}()

	log.Println("Бот запущен. Ожидание команд...")
	<-stop
	log.Println("Бот остановлен")
}
