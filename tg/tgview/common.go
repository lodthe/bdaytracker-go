package tgview

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

func cancelKeyboard() [][]telegram.KeyboardButton {
	return [][]telegram.KeyboardButton{
		{
			{
				Text: btn.Menu,
			},
		},
	}
}
