package db

import "time"

type Reminder struct {
	UserID         int64
	Text           string
	Category       string
	RemindAt       time.Time
	RepeatInterval string
}
