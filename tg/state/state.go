package state

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/lodthe/bdaytracker-go/models"
)

type State struct {
	State       ID // Conversation state
	StateBefore ID // State before handling current update

	TelegramID   int
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string

	CannotReceiveMessages bool

	LastNotificationAt time.Time

	VKID int

	Friends []models.Friend

	NewFriend models.Friend
}

func (s *State) Value() (driver.Value, error) {
	return json.Marshal(*s)
}

func (s *State) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("value cannot be converted to []byte")
	}

	return json.Unmarshal(b, s)
}
