package db

import tb "gopkg.in/telebot.v3"

func SaveUser(u *tb.User) error {
	_, err := DB.Exec(`
		INSERT INTO users (user_id, username, first_name, last_name, language)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			language = EXCLUDED.language
	`,
		u.ID,
		u.Username,
		u.FirstName,
		u.LastName,
		u.LanguageCode,
	)

	return err
}
