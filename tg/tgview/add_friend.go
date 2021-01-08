package tgview

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/models"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type AddFriend struct {
}

func (a AddFriend) AskName(s *tg.Session) {
	a.send(s, "Отправь имя нового 🧑‍\U0001F9B0 друга или 👩‍\U0001F9B0 подруги.")
}

func (a AddFriend) AskDate(s *tg.Session) {
	a.send(s, `Отправь дату рождения друга или подруги в следующем формате:

<code>ДД.ММ</code>

Например, 09.07 означает девятое июля.
`)
}

func (a AddFriend) FailedToParseDate(s *tg.Session) {
	a.send(s, `Не могу понять, что ты имеешь ввиду 😔
Сообщение должно соответствовать следующему формату:
<code>ДД.ММ</code>

Попробуй еще раз! 😉`)
}

func (a AddFriend) WrongNumberOfDays(s *tg.Session) {
	a.send(s, `❌ В этом месяце не может быть столько дней. Попробуй еще раз!😉`)
}

func (a AddFriend) Cancel(s *tg.Session) {
	_ = s.SendText(`Отменяю. 

Может, как-нибудь в следующий раз.`, Menu{}.Keyboard())
}

func (a AddFriend) Success(s *tg.Session, newFriend models.Friend) {
	keyboard := [][]telegram.InlineKeyboardButton{
		{
			callback.Button(btn.AddFriend, callback.AddFriend{}),
		},
		{
			callback.Button(btn.FriendList, callback.FriendList{}),
			callback.Button(btn.Menu, callback.OpenMenu{}),
		},
	}

	_ = s.SendText("👥", telegram.ReplyKeyboardRemove{})
	_ = s.SendText("<code>"+newFriend.Name+"</code> успешно добавлен(а) в список друзей!", keyboard)
}

func (a AddFriend) send(s *tg.Session, text string) {
	_ = s.SendText(text, [][]telegram.KeyboardButton{{
		{
			Text: btn.Cancel,
		},
	}})
}
