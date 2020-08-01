package handle

import (
	"sync"

	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/tg"
)

type UpdatesCollector struct {
	sessionLockers map[int]sync.Locker
	lock           sync.Locker
}

func NewUpdatesCollector() *UpdatesCollector {
	return &UpdatesCollector{
		sessionLockers: map[int]sync.Locker{},
		lock:           &sync.Mutex{},
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
			c.lock.Lock()
			sessionLocker, exists := c.sessionLockers[telegramID]
			if !exists {
				sessionLocker = &sync.Mutex{}
				c.sessionLockers[telegramID] = sessionLocker
			}
			c.lock.Unlock()

			sessionLocker.Lock()
			defer sessionLocker.Unlock()
			dispatchUpdate(&general, telegramID, update)
		}(userTelegramID, upd)
	}
}
