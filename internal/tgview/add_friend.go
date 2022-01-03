package tgview

import (
	friendship2 "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

type AddFriend struct {
}

func (a AddFriend) AskName(s *usersession.Session) {
	a.send(s, "Отправь имя нового 🧑‍\U0001F9B0 друга или 👩‍\U0001F9B0 подруги.")
}

func (a AddFriend) AskDate(s *usersession.Session) {
	a.send(s, `Отправь дату рождения друга или подруги в следующем формате:

<code>ДД.ММ</code>

Например, 09.07 означает девятое июля.
`)
}

func (a AddFriend) FailedToParseDate(s *usersession.Session) {
	a.send(s, `Не могу понять, что ты имеешь ввиду 😔
Сообщение должно соответствовать следующему формату:
<code>ДД.ММ</code>

Попробуй еще раз! 😉`)
}

func (a AddFriend) WrongNumberOfDays(s *usersession.Session) {
	a.send(s, `❌ В этом месяце не может быть столько дней. Попробуй еще раз!😉`)
}

func (a AddFriend) Cancel(s *usersession.Session) {
	_ = s.SendText(`Отменяю. 

Может, как-нибудь в следующий раз.`, Menu{}.Keyboard())
}

func (a AddFriend) Success(s *usersession.Session, newFriend friendship2.Friend) {
	keyboard := [][]telegram.InlineKeyboardButton{
		{
			tgcallback.Button(btn.AddFriend, tgcallback.AddFriend{}),
		},
		{
			tgcallback.Button(btn.FriendList, tgcallback.FriendList{}),
			tgcallback.Button(btn.Menu, tgcallback.OpenMenu{}),
		},
	}

	_ = s.SendText("👥", menuKeyboard())
	_ = s.SendText("<code>"+newFriend.Name+"</code> успешно добавлен(а) в список друзей!", keyboard)
}

func (a AddFriend) send(s *usersession.Session, text string) {
	_ = s.SendText(text, [][]telegram.KeyboardButton{{
		{
			Text: btn.Cancel,
		},
	}})
}
