package handle

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/lodthe/bdaytracker-go/models"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
	"github.com/lodthe/bdaytracker-go/tg/tgview/btn"
)

type AddFriendHandler struct {
}

func (h *AddFriendHandler) State() state.ID {
	return state.AddFriend
}

func (h *AddFriendHandler) Callback() interface{} {
	return callback.AddFriend{}
}

func (h *AddFriendHandler) HandleCallback(s *tg.Session, clb interface{}) {
	s.State.State = state.AddFriend
	s.State.NewFriend = models.Friend{}
	tgview.AddFriend{}.AskName(s)
}

func (h *AddFriendHandler) HandleMessage(s *tg.Session, msgText string) {
	switch {
	case msgText == btn.Cancel:
		h.cancel(s)

	case s.State.NewFriend.Name == "":
		h.handleName(s, msgText)

	case s.State.NewFriend.BDay == 0:
		h.handleDate(s, msgText)
	}
}

func (h *AddFriendHandler) cancel(s *tg.Session) {
	s.State.State = state.None
	tgview.AddFriend{}.Cancel(s)
}

func (h *AddFriendHandler) handleName(s *tg.Session, msgText string) {
	s.State.NewFriend.Name = msgText
	tgview.AddFriend{}.AskDate(s)
}

func (h *AddFriendHandler) handleDate(s *tg.Session, msgText string) {
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

	friend.UUID = fmt.Sprint(uuid.New())
	s.State.Friends = append(s.State.Friends, friend)
	s.State.State = state.None
	s.State.NewFriend = models.Friend{}

	tgview.AddFriend{}.Success(s, friend)
}
