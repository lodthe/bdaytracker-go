package tghandle

import (
	"reflect"
	"runtime/debug"

	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/petuhovskiy/telegram"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
)

func dispatchUpdate(s *usersession.Session, update telegram.Update) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"recovered":   r,
				"telegram_id": s.TelegramID,
				"stacktrace":  string(debug.Stack()),
				"update":      update,
			}).Error("recovered from panic")
		}
	}()

	if update.CallbackQuery != nil {
		clb := tgcallback.Unmarshal(update.CallbackQuery.Data)

		logrus.WithFields(logrus.Fields{
			"telegram_id": s.TelegramID,
			"type_name":   reflect.TypeOf(clb).Name(),
		}).Info("callback unpacked")
	}

	s.State.StateBefore = s.State.State
	s.State.CannotReceiveMessages = false // If we receive an update from the user, they can receive our messages

	updateUserInfo(s, update)

	s.AnswerOnLastCallback()
	activateHandler(s, update,
		&StartHandler{},
		&AddFriendHandler{},
		&FriendListHandler{},

		&RemoveFriendHandler{},
		&RemoveFriendApproveHandler{},
		&RemoveFriendCancelHandler{},
		&AddFromVKHandler{},
		&RemoveFromVKHandler{},

		&MenuHandler{},
	)

	err := s.SaveState()
	if err != nil {
		logrus.WithField("telegram_id", s.TelegramID).WithError(err).Error("failed to save the state")
	}
}

// updateUserInfo saves the user's name, username and language code.
func updateUserInfo(s *usersession.Session, update telegram.Update) {
	var user *telegram.User

	switch {
	case update.Message != nil && update.Message.From != nil:
		user = update.Message.From

	case update.CallbackQuery != nil && update.CallbackQuery.From != nil:
		user = update.CallbackQuery.From
	}

	if user == nil {
		return
	}

	s.State.Username = user.Username
	s.State.FirstName = user.FirstName
	s.State.LastName = user.LastName
	s.State.LanguageCode = user.LanguageCode
}
