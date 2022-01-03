package tghandle

import (
	friendship "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type RemoveFromVKHandler struct {
}

func (h *RemoveFromVKHandler) Callback() interface{} {
	return tgcallback.RemoveFromVK{}
}

func (h *RemoveFromVKHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	s.State.Friends = friendship.RemoveVKFriends(s.State.Friends)
	tgview.RemoveFromVK{}.Success(s)
}
