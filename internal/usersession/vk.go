package usersession

import (
	"github.com/lodthe/bdaytracker-go/internal/friendship"
)

func (s *Session) GetVKFriends(id int) ([]friendship.Friend, error) {
	return s.ctrl.vkCli.GetFriends(id)
}
