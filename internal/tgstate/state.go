package tgstate

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	friendship "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/pkg/errors"
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

	Friends []friendship.Friend

	NewFriend friendship.Friend
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
