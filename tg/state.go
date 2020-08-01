package tg

import (
	"github.com/jinzhu/gorm"
	"github.com/lodthe/bdaytracker-go/tg/state"
)

type State struct {
	gorm.Model
	TelegramID int
	State      state.State // Conversation state
}

func loadState(db *gorm.DB, telegramID int) (*State, error) {
	var state State
	err := db.Where(&State{
		TelegramID: telegramID,
	}).FirstOrCreate(&state, State{
		TelegramID: telegramID,
	}).Error

	return &state, err
}

func (s *State) Save(db *gorm.DB) error {
	return db.Save(s).Error
}
