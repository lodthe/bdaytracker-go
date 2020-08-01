package tg

type Session struct {
	General

	TelegramID int

	State *State
}

func NewSession(telegramID int, general *General) (*Session, error) {
	state, err := loadState(general.DB, telegramID)
	if err != nil {
		return nil, err
	}

	return &Session{
		General:    *general,
		TelegramID: telegramID,
		State:      state,
	}, nil
}

func (s *Session) SaveState() error {
	return s.State.Save(s.DB)
}
