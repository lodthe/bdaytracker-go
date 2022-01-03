package tgview

import (
	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type AddFromVK struct {
}

func (f AddFromVK) AskForID(s *usersession.Session) {
	text := `Отправь свой ID профиля ВКонтакте, чтобы я смог получить информацию о твоих друзья.

Узнать свой ID можно <a href="https://regvk.com/id/">здесь</a>.`

	s.SendText(text, cancelKeyboard())
}

func (f AddFromVK) IDIsNotANumber(s *usersession.Session) {
	s.SendText("ID может состоять только из цифр.", cancelKeyboard())
}

func (f AddFromVK) ProfileIsHidden(s *usersession.Session) {
	s.SendText("Чтобы добавить друзей из ВКонтакте, профиль должен быть открытым. После успешного добавления информации о друзьях ты можешь закрыть профиль.", cancelKeyboard())
}

func (f AddFromVK) Canceled(s *usersession.Session) {
	s.SendText("Добавление друзей из ВКонтакте отменено.", menuKeyboard())
}

func (f AddFromVK) Success(s *usersession.Session) {
	s.SendText("Информация о друзьях из ВКонтакте обновлена успешно!", menuKeyboard())
}
