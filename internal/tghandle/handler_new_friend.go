package tghandle

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	friendship2 "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/lodthe/bdaytracker-go/internal/usersession"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
	"github.com/lodthe/bdaytracker-go/internal/tgview/btn"
)

type AddFriendHandler struct {
}

func (h *AddFriendHandler) State() tgstate.ID {
	return tgstate.AddFriend
}

func (h *AddFriendHandler) Callback() interface{} {
	return tgcallback.AddFriend{}
}

func (h *AddFriendHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	s.State.State = tgstate.AddFriend
	s.State.NewFriend = friendship2.Friend{}
	tgview.AddFriend{}.AskName(s)
}

func (h *AddFriendHandler) HandleMessage(s *usersession.Session, msgText string) {
	switch {
	case msgText == btn.Cancel:
		h.cancel(s)

	case s.State.NewFriend.Name == "":
		h.handleName(s, msgText)

	case s.State.NewFriend.BDay == 0:
		h.handleDate(s, msgText)
	}
}

func (h *AddFriendHandler) cancel(s *usersession.Session) {
	s.State.State = tgstate.None
	tgview.AddFriend{}.Cancel(s)
}

func (h *AddFriendHandler) handleName(s *usersession.Session, msgText string) {
	s.State.NewFriend.Name = msgText
	tgview.AddFriend{}.AskDate(s)
}

func (h *AddFriendHandler) handleDate(s *usersession.Session, msgText string) {
	const numberOfMonths = 12
	var daysBefore = [...]int{ // It's took from the time package
		0,
		31,
		31 + 28,
		31 + 28 + 31,
		31 + 28 + 31 + 30,
		31 + 28 + 31 + 30 + 31,
		31 + 28 + 31 + 30 + 31 + 30,
		31 + 28 + 31 + 30 + 31 + 30 + 31,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
	}

	friend := s.State.NewFriend
	_, err := fmt.Sscanf(msgText, "%d.%d", &friend.BDay, &friend.BMonth)
	if err != nil {
		tgview.AddFriend{}.FailedToParseDate(s)
		return
	}

	if friend.BMonth < 1 || friend.BMonth > numberOfMonths {
		tgview.AddFriend{}.FailedToParseDate(s)
		return
	}

	daysInMonth := daysBefore[friend.BMonth] - daysBefore[friend.BMonth-1]
	if friend.BMonth == int(time.February) {
		daysInMonth++
	}
	if friend.BDay < 0 || friend.BDay > daysInMonth {
		tgview.AddFriend{}.WrongNumberOfDays(s)
		return
	}

	friend.UUID = uuid.New().String()
	s.State.Friends = append(s.State.Friends, friend)
	s.State.State = tgstate.None
	s.State.NewFriend = friendship2.Friend{}

	tgview.AddFriend{}.Success(s, friend)
}
