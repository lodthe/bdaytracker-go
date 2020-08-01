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
	var st State
	err := db.Where(&State{
		TelegramID: telegramID,
	}).FirstOrCreate(&st, State{
		TelegramID: telegramID,
	}).Error

	return &st, err
}

func (s *State) Save(db *gorm.DB) error {
	return db.Save(s).Error
}
