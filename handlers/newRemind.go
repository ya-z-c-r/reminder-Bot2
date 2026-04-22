package handlers

import (
	"log"
	"reminder-bot/db"
)

func newRemind(r db.Reminder) error {
	_, err := db.DB.Exec(`
		INSERT INTO reminders (user_id, text, category, remind_at, repeat_interval)
		VALUES ($1, $2, $3, $4, $5)
	`,
		r.UserID,
		r.Text,
		r.Category,
		r.RemindAt,
		r.RepeatInterval,
	)
	if err != nil {
		log.Fatal("ошибка записи напоминания в бд")
		return err
	}
	return nil
}
