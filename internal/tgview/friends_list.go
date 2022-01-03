package tgview

import (
	"github.com/petuhovskiy/telegram"

	friendship "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/usersession"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

const pageSize int = 15

type FriendList struct {
}

func (f FriendList) Send(s *usersession.Session, clb tgcallback.FriendList) {
	clb.Offset = minInt(clb.Offset, len(s.State.Friends)-1)
	clb.Offset = maxInt(clb.Offset, 0)

	sorted := friendship.SortFriends(s.State.Friends)

	var friends []friendship.Friend
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

func (f FriendList) keyboard(s *usersession.Session, clb tgcallback.FriendList) [][]telegram.InlineKeyboardButton {
	var prev interface{} = tgcallback.None{}
	var next interface{} = tgcallback.None{}
	if clb.Offset > 0 {
		prev = tgcallback.FriendList{
			Offset: maxInt(0, clb.Offset-pageSize),
		}
	}
	if clb.Offset+pageSize < len(s.State.Friends) {
		next = tgcallback.FriendList{
			Offset: clb.Offset + pageSize,
		}
	}

	// Insert pagination and delete_friend buttons if the Friends list is not empty
	var keyboard [][]telegram.InlineKeyboardButton
	if len(s.State.Friends) != 0 {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{
			tgcallback.Button(btn.Prev, prev),
			tgcallback.Button(btn.Next, next),
		})

		removeButtons := []telegram.InlineKeyboardButton{
			tgcallback.Button(btn.RemoveFriend, tgcallback.RemoveFriend{}),
		}
		if s.State.VKID != 0 {
			removeButtons = append(removeButtons, tgcallback.Button(btn.RemoveFromVK, tgcallback.RemoveFromVK{}))
		}

		keyboard = append(keyboard, removeButtons)
	}

	return append(keyboard,
		[]telegram.InlineKeyboardButton{
			tgcallback.Button(btn.AddFriend, tgcallback.AddFriend{}),
		},
		[]telegram.InlineKeyboardButton{
			tgcallback.Button(btn.Menu, tgcallback.OpenMenu{
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
