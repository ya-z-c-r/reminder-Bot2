package db

import "log"

func NewRemind(r Reminder) error {
	_, err := DB.Exec(`INSERT INTO reminders (text, remind_at, user_id)
	VALUES ($1, $2, $3)
	`,
		r.Text,
		r.RemindAt,
		r.UserID,
	)
	if err != nil {
		log.Print("ошибка при добавлении в базу даннный воспоминания", err)
	}
	return err
}
