package usersession

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/conf"
	"github.com/lodthe/bdaytracker-go/internal/tglimiter"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/vk"
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
