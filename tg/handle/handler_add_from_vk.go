package handle

import (
	"strconv"

	apiErrors "github.com/SevereCloud/vksdk/api/errors"

	"github.com/lodthe/bdaytracker-go/helpers"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type AddFromVKHandler struct {
}

func (h *AddFromVKHandler) State() state.ID {
	return state.ImportFromVK
}

func (h *AddFromVKHandler) Callback() interface{} {
	return callback.AddFromVK{}
}

func (h *AddFromVKHandler) HandleCallback(s *tg.Session, clb interface{}) {
	s.State.State = state.ImportFromVK
	tgview.AddFromVK{}.AskForID(s)
}

func (h *AddFromVKHandler) HandleMessage(s *tg.Session, msgText string) {
	switch msgText {
	case btn.Cancel:
		tgview.AddFromVK{}.Canceled(s)
		s.State.State = state.None

	default:
		h.handleVKID(s, msgText)
	}
}

func (h *AddFromVKHandler) handleVKID(s *tg.Session, vkID string) {
	id, err := strconv.Atoi(vkID)
	if err != nil {
		tgview.AddFromVK{}.IDIsNotANumber(s)
		return
	}

	s.State.VKID = id
	friendsToAdd, err := s.VKCli.GetFriends(id)
	if apiErrors.GetType(err) == apiErrors.PrivateProfile {
		tgview.AddFromVK{}.ProfileIsHidden(s)
		return
	}
	if err != nil {
		tgview.AddFromVK{}.AskForID(s)
		return
	}

	s.State.Friends = helpers.RemoveVKFriends(s.State.Friends)
	s.State.Friends = append(s.State.Friends, friendsToAdd...)
	s.State.State = state.None
	tgview.AddFromVK{}.Success(s)
}
