package tgview

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

func SendStartMessage(s *tg.Session) {
	s.SendText(`Привет! Я умею напоминать про 🎁 Дни Рождения твоих друзей.

Ты можешь добавить информацию о 🎁 Дне Рождения друга самостоятельно или импортировать даты рождения своих друзей из ВКонтакте.
        
Когда наступит чей-то День Рождения, я напомню тебе об этом!`, [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.Menu, callback.OpenMenu{}),
		},
	})
}
