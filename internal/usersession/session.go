package usersession

import (
	vk2 "github.com/lodthe/bdaytracker-go/internal/vk"
	"github.com/petuhovskiy/telegram"
	"github.com/pkg/errors"

	"github.com/lodthe/bdaytracker-go/internal/tglimiter"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
)

type Session struct {
	VKCli    *vk2.Client
	Bot      *telegram.Bot
	Executor *tglimiter.Executor

	TelegramID int
	LastUpdate *telegram.Update

	State *tgstate.State
}

func NewSession(vkCli *vk2.Client, bot *telegram.Bot, executor *tglimiter.Executor, repo tgstate.Repository, telegramID int, update *telegram.Update) (*Session, error) {
	st, err := repo.Get(telegramID)
	if err != nil {
		return nil, errors.Wrap(err, "state loading failed")
	}

	return &Session{
		VKCli:      vkCli,
		Bot:        bot,
		Executor:   executor,
		TelegramID: telegramID,
		LastUpdate: update,
		State:      st,
	}, nil
}
