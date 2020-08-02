package tgview

import (
	"fmt"

	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

func SendMenu(s *tg.Session) {
	text := fmt.Sprintf(`<b>%s</b>

<b>%s</b> — добавить вручную нового друга.
        
<b>%s</b> — обновить данные о друзьях из ВКонтакте.
        
<b>%s</b> — просмотреть список уже добавленных друзей.`, btn.Menu, btn.AddFriend, btn.AddFromVK, btn.FriendsList)
	keyboard := [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.AddFriend, callback.AddFriend{}),
			callback.Button(btn.AddFromVK, callback.AddFromVK{}),
		},
		{
			callback.Button(btn.FriendsList, callback.FriendsList{}),
		},
		{
			callback.Button(btn.Settings, callback.Settings{}),
		},
	}

	s.SendText(text, keyboard)
}
