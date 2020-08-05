package tgview

import (
	"github.com/lodthe/bdaytracker-go/tg"
)

type Start struct {
}

func (st Start) Send(s *tg.Session) {
	s.SendInlinePhoto(`Привет! Я умею напоминать про 🎁 Дни Рождения твоих друзей.

Ты можешь добавить информацию самостоятельно или импортировать из ВКонтакте.
        
Когда наступит чей-то День Рождения, я напомню тебе об этом!`, "greetings.png", nil)
}
