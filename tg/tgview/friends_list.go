package tgview

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/helpers"
	"github.com/lodthe/bdaytracker-go/models"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

const pageSize int = 15

type FriendList struct {
}

func (f FriendList) Send(s *tg.Session, clb callback.FriendList) {
	clb.Offset = minInt(clb.Offset, len(s.State.Friends)-1)
	clb.Offset = maxInt(clb.Offset, 0)

	sorted := helpers.SortFriends(s.State.Friends)

	var friends []models.Friend
	if len(sorted) != 0 {
		friends = sorted[clb.Offset:minInt(clb.Offset+pageSize, len(sorted))]
	}

	var text string
	for i := range friends {
		text += formatFriendWithIndex(friends[i], clb.Offset+i+1, len(s.State.Friends)) + "\n"
	}

	if text == "" {
		text = `Пока что друзей нет 😒

Ты можешь добавить их в ` + btn.Menu
	}

	s.SendEditText(text, f.keyboard(s, clb), true)
}

func (f FriendList) keyboard(s *tg.Session, clb callback.FriendList) [][]telegram.InlineKeyboardButton {
	var prev interface{} = callback.None{}
	var next interface{} = callback.None{}
	if clb.Offset > 0 {
		prev = callback.FriendList{
			Offset: maxInt(0, clb.Offset-pageSize),
		}
	}
	if clb.Offset+pageSize < len(s.State.Friends) {
		next = callback.FriendList{
			Offset: clb.Offset + pageSize,
		}
	}

	// Insert pagination and delete_friend buttons if the Friends list is not empty
	var keyboard [][]telegram.InlineKeyboardButton
	if len(s.State.Friends) != 0 {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{
			callback.Button(btn.Prev, prev),
			callback.Button(btn.Next, next),
		})

		removeButtons := []telegram.InlineKeyboardButton{
			callback.Button(btn.RemoveFriend, callback.RemoveFriend{}),
		}
		if s.State.VKID != 0 {
			removeButtons = append(removeButtons, callback.Button(btn.RemoveFromVK, callback.RemoveFromVK{}))
		}

		keyboard = append(keyboard, removeButtons)
	}

	return append(keyboard,
		[]telegram.InlineKeyboardButton{
			callback.Button(btn.AddFriend, callback.AddFriend{}),
		},
		[]telegram.InlineKeyboardButton{
			callback.Button(btn.Menu, callback.OpenMenu{
				Edit: true,
			},
			),
		})
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
