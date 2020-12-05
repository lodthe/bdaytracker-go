package tgview

import (
	"fmt"

	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type Menu struct {
}

func (m Menu) Send(s *tg.Session, edit bool) {
	text := fmt.Sprintf(`<b>%s</b>

<b>%s</b> — добавить вручную нового друга.
        
<b>%s</b> — обновить данные о друзьях из ВКонтакте.
        
<b>%s</b> — просмотреть список уже добавленных друзей.`, btn.Menu, btn.AddFriend, btn.AddFromVK, btn.FriendList)
	keyboard := [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.AddFriend, callback.AddFriend{}),
			callback.Button(btn.AddFromVK, callback.AddFromVK{}),
		},
		{
			callback.Button(btn.FriendList, callback.FriendList{}),
		},
		{
			callback.Button(btn.Settings, callback.Settings{}),
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
