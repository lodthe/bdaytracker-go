package handle

import (
	"runtime/debug"

	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/tg"
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

	s, err := tg.NewSession(sessionTelegramID, general)
	if err != nil {
		log.WithFields(log.Fields{
			"session_telegram_id": sessionTelegramID,
			"update":              update,
		}).WithError(err).Error("failed to create the session")
		return
	}

	activateHandler(s, update,
		&StartHandler{},
	)
}
