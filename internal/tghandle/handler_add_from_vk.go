package tghandle

import (
	"strconv"

	apiErrors "github.com/SevereCloud/vksdk/api/errors"

	"github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/usersession"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

type AddFromVKHandler struct {
}

func (h *AddFromVKHandler) State() tgstate.ID {
	return tgstate.ImportFromVK
}

func (h *AddFromVKHandler) Callback() interface{} {
	return tgcallback.AddFromVK{}
}

func (h *AddFromVKHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	s.State.State = tgstate.ImportFromVK
	tgview.AddFromVK{}.AskForID(s)
}

func (h *AddFromVKHandler) HandleMessage(s *usersession.Session, msgText string) {
	switch msgText {
	case btn.Cancel:
		tgview.AddFromVK{}.Canceled(s)
		s.State.State = tgstate.None

	default:
		h.handleVKID(s, msgText)
	}
}

func (h *AddFromVKHandler) handleVKID(s *usersession.Session, vkID string) {
	id, err := strconv.Atoi(vkID)
	if err != nil {
		tgview.AddFromVK{}.IDIsNotANumber(s)
		return
	}

	s.State.VKID = id
	friendsToAdd, err := s.GetVKFriends(id)
	if apiErrors.GetType(err) == apiErrors.PrivateProfile {
		tgview.AddFromVK{}.ProfileIsHidden(s)
		return
	}
	if err != nil {
		tgview.AddFromVK{}.AskForID(s)
		return
	}

	s.State.Friends = friendship.RemoveVKFriends(s.State.Friends)
	s.State.Friends = append(s.State.Friends, friendsToAdd...)
	s.State.State = tgstate.None
	tgview.AddFromVK{}.Success(s)
}
