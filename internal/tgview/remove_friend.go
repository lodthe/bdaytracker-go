package tgview

import (
	"fmt"

	friendship "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

type RemoveFriend struct {
}

func (f RemoveFriend) AskIndexOrName(s *usersession.Session) {
	s.SendText("Отправь полное имя друга или его номер из списка друзей.", cancelKeyboard())
}

func (f RemoveFriend) WrongIndex(s *usersession.Session) {
	s.SendText("Неправильный номер. Попробуй еще раз!", cancelKeyboard())
}

func (f RemoveFriend) WrongName(s *usersession.Session) {
	s.SendText("Не могу найти друга с таким именем. Имя должно быть таким же, как и в списке друзей.\n\nПопробуй еще раз!", cancelKeyboard())
}

func (f RemoveFriend) AskForApprove(s *usersession.Session, friend friendship.Friend) {
	text := fmt.Sprintf("Из списка друзей будет удалена запись\n%s", formatFriend(friend))
	s.SendText(text, [][]telegram.InlineKeyboardButton{
		{
			tgcallback.Button(btn.Approve, tgcallback.RemoveFriendApprove{
				UUID: friend.UUID,
				Name: friend.Name,
			}),
			tgcallback.Button(btn.Cancel, tgcallback.RemoveFriendCancel{
				Name: friend.Name,
			}),
		},
	})
}

func (f RemoveFriend) Approved(s *usersession.Session, clb tgcallback.RemoveFriendApprove) {
	text := fmt.Sprintf("<b>%s</b> удален(а) из списка друзей.", clb.Name)
	_ = s.DeleteLastMessage()
	_ = s.SendText(text, menuKeyboard())
}

func (f RemoveFriend) Canceled(s *usersession.Session, clb tgcallback.RemoveFriendCancel) {
	text := fmt.Sprintf("<b>%s</b> остается в списке друзей!", clb.Name)
	_ = s.DeleteLastMessage()
	_ = s.SendText(text, menuKeyboard())
}

func (f RemoveFriend) Cancel(s *usersession.Session) {
	_ = s.SendText("Отменено.", Menu{}.Keyboard())
}
