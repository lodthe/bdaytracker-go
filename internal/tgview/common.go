package tgview

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

func cancelKeyboard() [][]telegram.KeyboardButton {
	return [][]telegram.KeyboardButton{
		{
			{
				Text: btn.Cancel,
			},
		},
	}
}

func menuKeyboard() [][]telegram.KeyboardButton {
	return [][]telegram.KeyboardButton{
		{
			{
				Text: btn.Menu,
			},
		},
	}
}
