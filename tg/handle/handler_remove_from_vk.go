package handle

import (
	"github.com/lodthe/bdaytracker-go/helpers"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

type RemoveFromVKHandler struct {
}

func (h *RemoveFromVKHandler) Callback() interface{} {
	return callback.RemoveFromVK{}
}

func (h *RemoveFromVKHandler) HandleCallback(s *tg.Session, clb interface{}) {
	s.State.Friends = helpers.RemoveVKFriends(s.State.Friends)
	tgview.RemoveFromVK{}.Success(s)
}
