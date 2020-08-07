package tgview

import (
	"fmt"

	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/models"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type RemoveFriend struct {
}

func (f RemoveFriend) AskIndexOrName(s *tg.Session) {
	s.SendText("Отправь полное имя друга из списка друзей.", [][]telegram.KeyboardButton{
		{
			{
				Text: btn.Cancel,
			},
		},
	})
}

func (f RemoveFriend) WrongIndex(s *tg.Session) {
	s.SendText("Неправильный номер. Попробуй еще раз!", cancelKeyboard())
}

func (f RemoveFriend) WrongName(s *tg.Session) {
	s.SendText("Не могу найти друга с таким именем. Имя должно быть таким же, как и в списке друзей.\n\nПопробуй еще раз!", cancelKeyboard())
}

func (f RemoveFriend) AskForApprove(s *tg.Session, friend models.Friend) {
	text := fmt.Sprintf("Из списка друзей будет удалена запись\n%s", formatFriend(friend))
	s.SendText(text, [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.Approve, callback.RemoveFriendApprove{
				UUID: friend.UUID,
				Name: friend.Name,
			}),
			callback.Button(btn.Cancel, callback.RemoveFriendCancel{
				Name: friend.Name,
			}),
		},
	})
}

func (f RemoveFriend) Approved(s *tg.Session, clb callback.RemoveFriendApprove) {
	text := fmt.Sprintf("<b>%s</b> удален(а) из списка друзей.", clb.Name)
	s.SendEditText(text, nil, true)
}

func (f RemoveFriend) Cancelled(s *tg.Session, clb callback.RemoveFriendCancel) {
	text := fmt.Sprintf("<b>%s</b> остается в списке друзей!", clb.Name)
	s.SendEditText(text, nil, true)
}

func (f RemoveFriend) Cancel(s *tg.Session) {
	s.SendText("Отменено.", Menu{}.Keyboard())
}
