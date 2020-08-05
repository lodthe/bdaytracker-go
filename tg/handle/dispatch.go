package handle

import (
	"reflect"
	"runtime/debug"

	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/callback"
)

func dispatchUpdate(general *tg.General, sessionTelegramID int, update telegram.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"recovered":           r,
				"session_telegram_id": sessionTelegramID,
				"stacktrace":          string(debug.Stack()),
				"update":              update,
			}).Error("recovered from panic")
		}
	}()

	s, err := tg.NewSession(sessionTelegramID, general, &update)
	if err != nil {
		log.WithFields(log.Fields{
			"session_telegram_id": sessionTelegramID,
			"update":              update,
		}).WithError(err).Error("failed to create the session")
		return
	}

	if update.CallbackQuery != nil {
		clb := callback.Unmarshal(update.CallbackQuery.Data)
		log.WithFields(log.Fields{
			"telegram_id": sessionTelegramID,
			"type_name":   reflect.TypeOf(clb).Name(),
		}).Info("unpack a callback")
	}

	s.AnswerOnLastCallback()
	activateHandler(s, update,
		&StartHandler{},
		&AddFriendHandler{},
		&FriendsListHandler{},

		&MenuHandler{},
	)

	err = s.SaveState()
	if err != nil {
		log.WithField("telegram_id", sessionTelegramID).WithError(err).Error("failed to save the state")
	}
}
