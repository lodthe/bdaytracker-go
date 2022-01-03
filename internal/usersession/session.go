package usersession

import (
	"github.com/lodthe/bdaytracker-go/internal/conf"
	"github.com/lodthe/bdaytracker-go/internal/vk"
	"github.com/petuhovskiy/telegram"
	"github.com/pkg/errors"

	"github.com/lodthe/bdaytracker-go/internal/tglimiter"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
)

type controllers struct {
	cfg        *conf.Config
	tgBot      *telegram.Bot
	tgExecutor *tglimiter.Executor
	vkCli      *vk.Client

	repo tgstate.Repository
}

type Session struct {
	ctrl controllers

	TelegramID int
	LastUpdate *telegram.Update

	State *tgstate.State
}

func (s *Session) SaveState() error {
	return s.ctrl.repo.Save(s.State)
}

type Issuer struct {
	ctrl controllers
}

func NewIssuer(cfg *conf.Config, tgBot *telegram.Bot, tgExecutor *tglimiter.Executor, vkCli *vk.Client, repo tgstate.Repository) *Issuer {
	return &Issuer{
		ctrl: controllers{
			cfg:        cfg,
			tgBot:      tgBot,
			tgExecutor: tgExecutor,
			vkCli:      vkCli,
			repo:       repo,
		},
	}
}

func (s *Issuer) Issue(telegramID int, update *telegram.Update) (*Session, error) {
	st, err := s.ctrl.repo.Get(telegramID)
	if err != nil {
		return nil, errors.Wrap(err, "state loading failed")
	}

	return &Session{
		ctrl:       s.ctrl,
		TelegramID: telegramID,
		LastUpdate: update,
		State:      st,
	}, nil
}
