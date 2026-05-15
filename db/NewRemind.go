package db

import "log"

func NewRemind(r Reminder) error {
	_, err := DB.Exec(`INSERT INTO reminders (text, remind_at, user_id, repeat_interval, done)
	VALUES ($1, $2, $3, $4, $5)
	`,
		r.Text,
		r.RemindAt,
		r.UserID,
		r.RepeatInterval,
		r.Done,
	)
	if err != nil {
		log.Print("ошибка при добавлении в базу даннный воспоминания", err)
	}
	return err
}
