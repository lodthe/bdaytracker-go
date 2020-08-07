package handle

import (
	"sort"
	"strconv"
	"strings"

	"github.com/lodthe/bdaytracker-go/models"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type RemoveFriendHandler struct {
}

func (h *RemoveFriendHandler) State() state.ID {
	return state.RemoveFriend
}

func (h *RemoveFriendHandler) Callback() interface{} {
	return callback.RemoveFriend{}
}

func (h *RemoveFriendHandler) HandleCallback(s *tg.Session, clb interface{}) {
	s.State.State = state.RemoveFriend
	tgview.RemoveFriend{}.AskIndexOrName(s)
}

func (h *RemoveFriendHandler) HandleMessage(s *tg.Session, msgText string) {
	switch msgText {
	case btn.Cancel:
		tgview.RemoveFriend{}.Cancel(s)
		s.State.State = state.None

	default:
		h.handleIndexOrName(s, msgText)
	}
}

func (h *RemoveFriendHandler) findByIndex(s *tg.Session, index int, friends []models.Friend) (friend models.Friend, found bool) {
	if index < 1 || index > len(friends) {
		tgview.RemoveFriend{}.WrongIndex(s)
		return models.Friend{}, false
	}
	return friends[index-1], true
}

func (h *RemoveFriendHandler) findByName(s *tg.Session, name string, friends []models.Friend) (friend models.Friend, found bool) {
	for i := range friends {
		if strings.EqualFold(friends[i].Name, name) {
			return friends[i], true
		}
	}
	tgview.RemoveFriend{}.WrongName(s)
	return models.Friend{}, false
}

func (h *RemoveFriendHandler) handleIndexOrName(s *tg.Session, msgText string) {
	friends := &tgview.FriendsArray{Friends: s.State.Friends}
	sort.Sort(friends)

	var friend models.Friend
	var found bool

	index, err := strconv.Atoi(msgText)
	if err != nil {
		friend, found = h.findByName(s, msgText, friends.Friends)
	} else {
		friend, found = h.findByIndex(s, index, friends.Friends)
	}

	if !found {
		return
	}

	tgview.RemoveFriend{}.AskForApprove(s, friend)
	s.State.State = state.None
}

type RemoveFriendApproveHandler struct {
}

func (h *RemoveFriendApproveHandler) Callback() interface{} {
	return callback.RemoveFriendApprove{}
}

func (h *RemoveFriendApproveHandler) HandleCallback(s *tg.Session, clb interface{}) {
	data := clb.(callback.RemoveFriendApprove)

	removeFriend := func(index int) {
		s.State.Friends = append(s.State.Friends[:index], s.State.Friends[index+1:]...)
	}

	for i := range s.State.Friends {
		if s.State.Friends[i].UUID == data.UUID {
			removeFriend(i)
			break
		}
	}

	tgview.RemoveFriend{}.Approved(s, clb.(callback.RemoveFriendApprove))
}

type RemoveFriendCancelHandler struct {
}

func (h *RemoveFriendCancelHandler) Callback() interface{} {
	return callback.RemoveFriendCancel{}
}

func (h *RemoveFriendCancelHandler) HandleCallback(s *tg.Session, clb interface{}) {
	tgview.RemoveFriend{}.Cancelled(s, clb.(callback.RemoveFriendCancel))
}
