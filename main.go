package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mustToken() string {
	token := flag.String("token",
		"",
		"токен телеграмм бота",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("токена не обнаружено")
	}

	return *token
}

func main() {
	botToken := mustToken()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Включаем режим отладки (опционально)
	bot.Debug = true

	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 2000

	updates := bot.GetUpdatesChan(u)

	// Обработка сигнала для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ping"),
			tgbotapi.NewKeyboardButton("Добавить напоминание"),
			tgbotapi.NewKeyboardButton("показать расписание на день"),
			tgbotapi.NewKeyboardButton("изменить напоминание"),
			tgbotapi.NewKeyboardButton("удалить напоминание"),
		),
	)

	go func() {
		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				switch update.Message.Command() {
				case "ping":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
					}
				case "start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давайте начнём")
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				}
			} else if update.Message != nil && !update.Message.IsCommand() {
				switch update.Message.Text {
				case "ping":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
					bot.Send(msg)
				}
			}
		}
	}()

	log.Println("Бот запущен. Ожидание команд...")
	<-stop
	log.Println("Бот остановлен")
}
