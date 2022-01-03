package tghandle

import (
	"github.com/petuhovskiy/telegram"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type UpdatesCollector struct {
	issuer         *usersession.Issuer
	sessionStorage *usersession.Storage
}

func NewUpdatesCollector(issuer *usersession.Issuer, storage *usersession.Storage) *UpdatesCollector {
	return &UpdatesCollector{
		issuer:         issuer,
		sessionStorage: storage,
	}
}

func (c *UpdatesCollector) Start(updates <-chan telegram.Update) {
	for upd := range updates {
		logrus.WithField("update", upd).Debug("received a new update")

		var userTelegramID int
		switch {
		case upd.Message != nil:
			userTelegramID = upd.Message.From.ID
		case upd.CallbackQuery != nil:
			userTelegramID = upd.CallbackQuery.From.ID
		default:
			continue
		}

		// For one session no more than 1 dispatcher can be in process at one moment
		go func(telegramID int, update telegram.Update) {
			sessionLocker := c.sessionStorage.AcquireLock(telegramID)
			sessionLocker.Lock()
			defer sessionLocker.Unlock()

			s, err := c.issuer.Issue(telegramID, &update)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"telegram_id": telegramID,
					"update":      update,
				}).WithError(err).Error("failed to issue a session")

				return
			}

			dispatchUpdate(s, update)
		}(userTelegramID, upd)
	}
}
