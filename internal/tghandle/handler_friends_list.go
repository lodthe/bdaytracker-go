package tghandle

import (
	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type FriendListHandler struct {
}

func (h *FriendListHandler) Callback() interface{} {
	return tgcallback.FriendList{}
}

func (h *FriendListHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	cdata := clb.(tgcallback.FriendList)
	tgview.FriendList{}.Send(s, cdata)
}
