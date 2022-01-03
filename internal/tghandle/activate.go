package tghandle

import (
	"reflect"

	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/petuhovskiy/telegram"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
)

type methodState interface {
	State() tgstate.ID // Returns the required state for handling
}

type methodCallback interface {
	Callback() interface{} // Returns a callback with the required type for handling
}

type methodCanHandle interface {
	CanHandle(s *usersession.Session, msg *telegram.Message, clb *telegram.CallbackQuery) bool // Means can the handler handle this update
}

type methodHandleMessage interface {
	HandleMessage(s *usersession.Session, msgText string)
}

type methodHandleCallback interface {
	HandleCallback(s *usersession.Session, clb interface{})
}

// activateHandlers goes through the given list of handlers (the same order as they are given)
// and stops when the current handler can handle the given update.
// Then it handles the update with the found handler, and search stops.
func activateHandler(s *usersession.Session, update telegram.Update, handlers ...interface{}) {
	activate := func(handler interface{}) {
		logger := logrus.WithFields(logrus.Fields{
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
				handler.HandleCallback(s, tgcallback.Unmarshal(update.CallbackQuery.Data))
			}

		default:
			logger.Error("the update can be handled, but a callback method is not provided")
		}
	}

	// ID-triggered conditions are more valuable than callback- or canHandle-triggered conditions.

	for i := range handlers {
		handler, ok := handlers[i].(methodState)
		if !ok || handler.State() != s.State.StateBefore {
			continue
		}
		activate(handlers[i])
		return
	}

	for i := range handlers {
		var canHandle bool

		handlerByCallback, ok := handlers[i].(methodCallback)
		if !canHandle && ok && update.CallbackQuery != nil {
			canHandle = reflect.TypeOf(tgcallback.Unmarshal(update.CallbackQuery.Data)) == reflect.TypeOf(handlerByCallback.Callback())
		}

		handlerByCanHandle, ok := handlers[i].(methodCanHandle)
		if !canHandle && ok {
			canHandle = handlerByCanHandle.CanHandle(s, update.Message, update.CallbackQuery)
		}

		if !canHandle {
			continue
		}

		activate(handlers[i])
		return
	}
}
