package tgview

import (
	"fmt"

	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/usersession"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

type Menu struct {
}

func (m Menu) Send(s *usersession.Session, edit bool) {
	text := fmt.Sprintf(`<b>%s</b>

<b>%s</b> — добавить вручную нового друга.
        
<b>%s</b> — обновить данные о друзьях из ВКонтакте.
        
<b>%s</b> — просмотреть список уже добавленных друзей.`, btn.Menu, btn.AddFriend, btn.AddFromVK, btn.FriendList)
	keyboard := [][]telegram.InlineKeyboardButton{
		{
			tgcallback.Button(btn.AddFriend, tgcallback.AddFriend{}),
			tgcallback.Button(btn.AddFromVK, tgcallback.AddFromVK{}),
		},
		{
			tgcallback.Button(btn.FriendList, tgcallback.FriendList{}),
		},
		{
			tgcallback.Button(btn.Settings, tgcallback.Settings{}),
		},
	}

	s.SendEditText(text, keyboard, edit)
}

func (m Menu) Keyboard() [][]telegram.KeyboardButton {
	return [][]telegram.KeyboardButton{
		{
			{
				Text: btn.Menu,
			},
		},
	}
}
