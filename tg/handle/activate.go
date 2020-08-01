package handle

import (
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"
)

type methodState interface {
	State() state.State // Returns the required state for handling
}

type methodCanHandle interface {
	CanHandle(s *tg.Session, msg *telegram.Message, clb *telegram.CallbackQuery) bool // Means can the handler handle this update
}

type methodHandleMessage interface {
	HandleMessage(s *tg.Session, msg *telegram.Message)
}

type methodHandleCallback interface {
	HandleCallback(s *tg.Session, clb *telegram.CallbackQuery)
}

// activateHandlers goes through the given list of handlers (the same order as they are given)
// and stops when the current handler can handle the given update.
// Then it handles the update with the found handler, and search stops.
func activateHandler(s *tg.Session, update telegram.Update, handlers ...interface{}) {
	for i := range handlers {
		var canHandle bool

		logger := log.WithFields(log.Fields{
			"update":      update,
			"telegram_id": s.TelegramID,
			"handler":     handlers[i],
		})

		switch handlers[i].(type) {
		case methodState:
			canHandle = s.State.State == handlers[i].(methodState).State()

		case methodCanHandle:
			canHandle = handlers[i].(methodCanHandle).CanHandle(s, update.Message, update.CallbackQuery)
		}

		if !canHandle {
			break
		}

		switch {
		case update.Message != nil:
			handler, ok := handlers[i].(methodHandleMessage)
			if !ok {
				logger.Error("missed HandleMessage method")
			} else {
				handler.HandleMessage(s, update.Message)
			}

		case update.CallbackQuery != nil:
			handler, ok := handlers[i].(methodHandleCallback)
			if !ok {
				logger.Error("missed HandleCallback method")
			} else {
				handler.HandleCallback(s, update.CallbackQuery)
			}
		}
	}
}
