package state

type State string

const (
	StateNone    State = ""
	StateAddText State = "add_text"
	StateAddTime State = "add_time"
)

type UserFlow struct {
	State State
	Text  string
}

var Flows = make(map[int64]*UserFlow)
