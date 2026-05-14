package ui

import tb "gopkg.in/telebot.v3"

var (
	MainMenu *tb.ReplyMarkup

	BtnAdd       tb.Btn
	BtnAddRepeat tb.Btn
)

func init() {
	MainMenu = &tb.ReplyMarkup{
		ResizeKeyboard: true,
	}

	BtnAdd = MainMenu.Text(
		"Добавить одноразовое напоминание",
	)

	BtnAddRepeat = MainMenu.Text(
		"Добавить повторяющееся напоминание",
	)

	MainMenu.Reply(
		MainMenu.Row(BtnAdd),
		MainMenu.Row(BtnAddRepeat),
	)
}
