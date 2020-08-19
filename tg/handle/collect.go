package handle

import (
	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/sessionstorage"
)

type UpdatesCollector struct {
	sessionStorage *sessionstorage.Storage
}

func NewUpdatesCollector(storage *sessionstorage.Storage) *UpdatesCollector {
	return &UpdatesCollector{
		sessionStorage: storage,
	}
}

func (c *UpdatesCollector) Start(general tg.General, updates <-chan telegram.Update) {
	for upd := range updates {
		log.WithField("update", upd).Debug("received a new update")

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
			dispatchUpdate(&general, telegramID, update)
		}(userTelegramID, upd)
	}
}
