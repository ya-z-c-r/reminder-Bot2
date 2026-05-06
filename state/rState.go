package state

import "time"

type State string

const (
	StateNone              State = ""
	StateAddText           State = "add_text"
	StateAddTime           State = "add_time"
	StateAddRepeatInterval State = "add_repeat_interval"
)

type UserFlow struct {
	State          State
	Text           string
	RemindAt       time.Time
	RepeatInterval string
}

var Flows = make(map[int64]*UserFlow)
