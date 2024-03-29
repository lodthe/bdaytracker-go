package tghandle

import (
	"github.com/petuhovskiy/telegram"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type SessionIssuer interface {
	Issue(telegramID int, update *telegram.Update) (s *usersession.Session, release func(), err error)
}

type UpdatesCollector struct {
	issuer SessionIssuer
}

func NewUpdatesCollector(issuer SessionIssuer) *UpdatesCollector {
	return &UpdatesCollector{
		issuer: issuer,
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
			s, release, err := c.issuer.Issue(telegramID, &update)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"telegram_id": telegramID,
					"update":      update,
				}).WithError(err).Error("failed to issue a session")

				return
			}
			defer release()

			dispatchUpdate(s, update)
		}(userTelegramID, upd)
	}
}
