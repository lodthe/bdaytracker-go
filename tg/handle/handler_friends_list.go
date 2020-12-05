package handle

import (
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

type FriendListHandler struct {
}

func (h *FriendListHandler) Callback() interface{} {
	return callback.FriendList{}
}

func (h *FriendListHandler) HandleCallback(s *tg.Session, clb interface{}) {
	cdata := clb.(callback.FriendList)
	tgview.FriendList{}.Send(s, cdata)
}
