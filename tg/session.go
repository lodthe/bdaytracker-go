package tg

import (
	"github.com/petuhovskiy/telegram"
)

type Session struct {
	General

	TelegramID int
	LastUpdate *telegram.Update

	State *State
}

func NewSession(telegramID int, general *General, update *telegram.Update) (*Session, error) {
	state, err := loadState(general.DB, telegramID)
	if err != nil {
		return nil, err
	}

	return &Session{
		General:    *general,
		TelegramID: telegramID,
		LastUpdate: update,
		State:      state,
	}, nil
}

func (s *Session) SaveState() error {
	return s.State.Save(s.DB)
}
