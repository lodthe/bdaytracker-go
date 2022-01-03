package tgview

import (
	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

type RemoveFromVK struct {
}

func (f RemoveFromVK) Success(s *usersession.Session) {
	_ = s.SendEditText("✅ Информация о друзьях из ВКонтакте удалена!", [][]telegram.InlineKeyboardButton{
		{
			tgcallback.Button(btn.FriendList, tgcallback.FriendList{}),
		},
		{
			tgcallback.Button(btn.Menu, tgcallback.OpenMenu{}),
		},
	}, true)
}
