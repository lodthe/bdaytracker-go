package handle

import (
	"strings"
	"time"

	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

const delayAfterStartMessage = time.Second

type StartHandler struct {
}

func (h *StartHandler) CanHandle(s *tg.Session, msg *telegram.Message, clb *telegram.CallbackQuery) bool {
	return msg != nil && strings.HasPrefix(msg.Text, "/start")
}

func (h *StartHandler) HandleMessage(s *tg.Session, msgText string) {
	tgview.Start{}.Send(s)
	time.Sleep(delayAfterStartMessage)
	tgview.Menu{}.Send(s, false)
}
