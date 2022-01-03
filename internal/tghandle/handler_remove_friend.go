package tghandle

import (
	"strconv"
	"strings"

	friendship2 "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type RemoveFriendHandler struct {
}

func (h *RemoveFriendHandler) State() tgstate.ID {
	return tgstate.RemoveFriend
}

func (h *RemoveFriendHandler) Callback() interface{} {
	return tgcallback.RemoveFriend{}
}

func (h *RemoveFriendHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	s.State.State = tgstate.RemoveFriend
	tgview.RemoveFriend{}.AskIndexOrName(s)
}

func (h *RemoveFriendHandler) HandleMessage(s *usersession.Session, msgText string) {
	switch msgText {
	case btn.Cancel:
		tgview.RemoveFriend{}.Cancel(s)
		s.State.State = tgstate.None

	default:
		h.handleIndexOrName(s, msgText)
	}
}

func (h *RemoveFriendHandler) findByIndex(s *usersession.Session, index int, friends []friendship2.Friend) (friend friendship2.Friend, found bool) {
	if index < 1 || index > len(friends) {
		tgview.RemoveFriend{}.WrongIndex(s)
		return friendship2.Friend{}, false
	}
	return friends[index-1], true
}

func (h *RemoveFriendHandler) findByName(s *usersession.Session, name string, friends []friendship2.Friend) (friend friendship2.Friend, found bool) {
	for i := range friends {
		if strings.EqualFold(friends[i].Name, name) {
			return friends[i], true
		}
	}
	tgview.RemoveFriend{}.WrongName(s)
	return friendship2.Friend{}, false
}

func (h *RemoveFriendHandler) handleIndexOrName(s *usersession.Session, msgText string) {
	sorted := friendship2.SortFriends(s.State.Friends)

	var friend friendship2.Friend
	var found bool

	index, err := strconv.Atoi(msgText)
	if err != nil {
		friend, found = h.findByName(s, msgText, sorted)
	} else {
		friend, found = h.findByIndex(s, index, sorted)
	}

	if !found {
		return
	}

	tgview.RemoveFriend{}.AskForApprove(s, friend)
	s.State.State = tgstate.None
}

type RemoveFriendApproveHandler struct {
}

func (h *RemoveFriendApproveHandler) Callback() interface{} {
	return tgcallback.RemoveFriendApprove{}
}

func (h *RemoveFriendApproveHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	data := clb.(tgcallback.RemoveFriendApprove)

	removeFriend := func(index int) {
		s.State.Friends = append(s.State.Friends[:index], s.State.Friends[index+1:]...)
	}

	for i := range s.State.Friends {
		if s.State.Friends[i].UUID == data.UUID {
			removeFriend(i)
			break
		}
	}

	tgview.RemoveFriend{}.Approved(s, clb.(tgcallback.RemoveFriendApprove))
}

type RemoveFriendCancelHandler struct {
}

func (h *RemoveFriendCancelHandler) Callback() interface{} {
	return tgcallback.RemoveFriendCancel{}
}

func (h *RemoveFriendCancelHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	tgview.RemoveFriend{}.Canceled(s, clb.(tgcallback.RemoveFriendCancel))
}
