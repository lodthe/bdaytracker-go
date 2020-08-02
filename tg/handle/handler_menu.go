package handle

import (
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

type MenuHandler struct {
}

func (h *MenuHandler) Callback() interface{} {
	return callback.OpenMenu{}
}

func (h *MenuHandler) HandleCallback(s *tg.Session, clb interface{}) {
	tgview.SendMenu(s)
}
