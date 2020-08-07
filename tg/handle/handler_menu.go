package handle

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

type MenuHandler struct {
}

func (h *MenuHandler) CanHandle(s *tg.Session, msg *telegram.Message, clb *telegram.CallbackQuery) bool {
	return s.State.StateBefore == state.None && msg != nil
}

func (h *MenuHandler) Callback() interface{} {
	return callback.OpenMenu{}
}

func (h *MenuHandler) HandleCallback(s *tg.Session, clb interface{}) {
	cdata := clb.(callback.OpenMenu)
	tgview.Menu{}.Send(s, cdata.Edit)
}

func (h *MenuHandler) HandleMessage(s *tg.Session, msgText string) {
	tgview.Menu{}.Send(s, false)
}
