package usersession

import (
	"github.com/lodthe/bdaytracker-go/internal/conf"
	"github.com/lodthe/bdaytracker-go/internal/tglimiter"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/vk"
	"github.com/petuhovskiy/telegram"
	"github.com/pkg/errors"
)

type Issuer struct {
	storage *storage
	ctrl    controllers
}

func NewIssuer(cfg *conf.Config, tgBot *telegram.Bot, tgExecutor *tglimiter.Executor, vkCli *vk.Client, repo tgstate.Repository) *Issuer {
	return &Issuer{
		storage: newStorage(),
		ctrl: controllers{
			cfg:        cfg,
			tgBot:      tgBot,
			tgExecutor: tgExecutor,
			vkCli:      vkCli,
			repo:       repo,
		},
	}
}

// Issue creates a new session and returns it.
// Also, it locks the session and returns the release function. Hence, you cannot work with a session simultaneously
// from difference threads.
// When you stop using the session, you must call release().
func (i *Issuer) Issue(telegramID int, update *telegram.Update) (s *Session, release func(), err error) {
	lock := i.storage.acquireLock(telegramID)
	lock.Lock()
	release = func() {
		lock.Unlock()
	}

	st, err := i.ctrl.repo.Get(telegramID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "state loading failed")
	}

	return &Session{
		ctrl:       i.ctrl,
		TelegramID: telegramID,
		LastUpdate: update,
		State:      st,
	}, release, nil
}
