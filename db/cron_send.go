package db

import (
	"log"

	"github.com/robfig/cron/v3"
	tb "gopkg.in/telebot.v3"
)

var (
	CronScheduler   *cron.Cron
	NewReminderChan = make(chan Reminder)
)

func StartCronWorker(bot *tb.Bot) {
	CronScheduler = cron.New()

	err := loadExistingReminders(bot)
	if err != nil {
		log.Println("cron load error:", err)
		return
	}

	CronScheduler.Start()

	log.Println("cron worker started")

	// слушаем новые напоминания
	for reminder := range NewReminderChan {
		AddCronReminder(bot, reminder)
	}
}

func loadExistingReminders(bot *tb.Bot) error {
	rows, err := DB.Query(`
		SELECT id, user_id, text, repeat_interval
		FROM reminders
		WHERE repeat_interval IS NOT NULL
		AND repeat_interval != ''
	`)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var reminder Reminder

		err := rows.Scan(
			&reminder.ID,
			&reminder.UserID,
			&reminder.Text,
			&reminder.RepeatInterval,
		)
		if err != nil {
			log.Println("scan error:", err)
			continue
		}

		AddCronReminder(bot, reminder)
	}

	return nil
}

func AddCronReminder(
	bot *tb.Bot,
	reminder Reminder,
) {
	// защита от closure бага
	r := reminder

	_, err := CronScheduler.AddFunc(
		r.RepeatInterval,
		func() {
			// log.Printf(
			// 	"cron triggered: user=%d text=%s",
			// 	r.UserID,
			// 	r.Text,
			// )

			_, err := bot.Send(
				&tb.User{ID: r.UserID},
				r.Text,
			)

			if err != nil {
				log.Println(
					"cron send error:",
					err,
				)
			}
		},
	)

	if err != nil {
		log.Println(
			"cron add error:",
			err,
			"interval:",
			r.RepeatInterval,
		)
		return
	}

	// log.Printf(
	// 	"cron added id=%d reminder=%d schedule=%s",
	// 	id,
	// 	r.ID,
	// 	r.RepeatInterval,
	// )
}
