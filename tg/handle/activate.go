package handle

import (
	"reflect"

	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/state"
)

type methodState interface {
	State() state.ID // Returns the required state for handling
}

type methodCallback interface {
	Callback() interface{} // Returns a callback with the required type for handling
}

type methodCanHandle interface {
	CanHandle(s *tg.Session, msg *telegram.Message, clb *telegram.CallbackQuery) bool // Means can the handler handle this update
}

type methodHandleMessage interface {
	HandleMessage(s *tg.Session, msgText string)
}

type methodHandleCallback interface {
	HandleCallback(s *tg.Session, clb interface{})
}

// activateHandlers goes through the given list of handlers (the same order as they are given)
// and stops when the current handler can handle the given update.
// Then it handles the update with the found handler, and search stops.
func activateHandler(s *tg.Session, update telegram.Update, handlers ...interface{}) {
	activate := func(handler interface{}) {
		logger := log.WithFields(log.Fields{
			"update":      update,
			"telegram_id": s.TelegramID,
			"handler":     handler,
		})

		switch {
		case update.Message != nil:
			handler, ok := handler.(methodHandleMessage)
			if !ok {
				logger.Error("missed HandleMessage method")
			} else {
				handler.HandleMessage(s, update.Message.Text)
			}

		case update.CallbackQuery != nil:
			handler, ok := handler.(methodHandleCallback)
			if !ok {
				logger.Error("missed HandleCallback method")
			} else {
				handler.HandleCallback(s, update.CallbackQuery.Data)
			}

		default:
			logger.Error("the update can be handled, but a callback method is not provided")
		}
	}

	// ID-triggered conditions are more valuable than callback- or canHandle-triggered conditions.

	for i := range handlers {
		handler, ok := handlers[i].(methodState)
		if !ok || handler.State() != s.State.State {
			continue
		}
		activate(handlers[i])
	}

	for i := range handlers {
		var canHandle bool

		switch handler := handlers[i].(type) {
		case methodCallback:
			if update.CallbackQuery != nil {
				canHandle = reflect.TypeOf(callback.Unmarshal(update.CallbackQuery.Data)) == reflect.TypeOf(handler.Callback())
			}

		case methodCanHandle:
			canHandle = handler.CanHandle(s, update.Message, update.CallbackQuery)
		}

		if !canHandle {
			continue
		}

		activate(handlers[i])
		return
	}
}
