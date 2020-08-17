package tgview

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type RemoveFromVK struct {
}

func (f RemoveFromVK) Success(s *tg.Session) {
	s.SendText("✅ Информация о друзьях из ВКонтакте удалена!", [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.FriendsList, callback.FriendsList{}),
		},
		{
			callback.Button(btn.Menu, callback.OpenMenu{}),
		},
	})
}
