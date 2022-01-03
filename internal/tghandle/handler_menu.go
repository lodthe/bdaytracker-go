package tghandle

import (
	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
)

type MenuHandler struct {
}

func (h *MenuHandler) CanHandle(s *usersession.Session, msg *telegram.Message, clb *telegram.CallbackQuery) bool {
	return s.State.StateBefore == tgstate.None && msg != nil
}

func (h *MenuHandler) Callback() interface{} {
	return tgcallback.OpenMenu{}
}

func (h *MenuHandler) HandleCallback(s *usersession.Session, clb interface{}) {
	cdata := clb.(tgcallback.OpenMenu)
	tgview.Menu{}.Send(s, cdata.Edit)
}

func (h *MenuHandler) HandleMessage(s *usersession.Session, msgText string) {
	tgview.Menu{}.Send(s, false)
}
