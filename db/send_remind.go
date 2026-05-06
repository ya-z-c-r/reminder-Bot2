package db

import (
	"log"
	"time"

	tb "gopkg.in/telebot.v3"
)

func GetDueReminders() ([]Reminder, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, text, remind_at
		FROM reminders
		WHERE done = false AND remind_at <= NOW()
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Reminder

	for rows.Next() {
		var r Reminder
		err := rows.Scan(&r.ID, &r.UserID, &r.Text, &r.RemindAt)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func MarkDone(id int) error {
	_, err := DB.Exec(`
		UPDATE reminders
		SET done = true
		WHERE id = $1
	`, id)

	return err
}

func StartReminderWorker(bot *tb.Bot) {
	for {
		reminders, err := GetDueReminders()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute)
			continue
		}

		for _, r := range reminders {

			if r.Text != "" {
				_, err := bot.Send(&tb.User{ID: r.UserID}, r.Text)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			err = MarkDone(r.ID)
			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(60 * time.Second)
	}
}
