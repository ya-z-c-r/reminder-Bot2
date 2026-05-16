package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	var err error

	connStr := "host=postgres port=5432 user=postgres password=password dbname=reminders sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS reminders (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		text TEXT NOT NULL,
		remind_at TIMESTAMP NOT NULL,
		repeat_interval TEXT,
		done BOOLEAN DEFAULT false
	);
	`

	_, err = DB.Exec(query)

	return err
}
