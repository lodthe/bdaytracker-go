package handle

import (
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

type FriendsListHandler struct {
}

func (h *FriendsListHandler) Callback() interface{} {
	return callback.FriendsList{}
}

func (h *FriendsListHandler) HandleCallback(s *tg.Session, clb interface{}) {
	cdata := clb.(callback.FriendsList)
	tgview.FriendsList{}.Send(s, cdata)
}
