package tg

import (
	"github.com/petuhovskiy/telegram"
	"github.com/pkg/errors"

	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/lodthe/bdaytracker-go/tg/tglimiter"
	"github.com/lodthe/bdaytracker-go/vk"
)

type Session struct {
	VKCli    *vk.Client
	Bot      *telegram.Bot
	Executor *tglimiter.Executor

	TelegramID int
	LastUpdate *telegram.Update

	State *state.State
}

func NewSession(vkCli *vk.Client, bot *telegram.Bot, executor *tglimiter.Executor, repo state.Repository, telegramID int, update *telegram.Update) (*Session, error) {
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
