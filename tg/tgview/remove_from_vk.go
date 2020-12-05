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
	_ = s.SendEditText("✅ Информация о друзьях из ВКонтакте удалена!", [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.FriendList, callback.FriendList{}),
		},
		{
			callback.Button(btn.Menu, callback.OpenMenu{}),
		},
	}, true)
}
