package tg

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg/state"

	"github.com/lodthe/bdaytracker-go/conf"
	"github.com/lodthe/bdaytracker-go/tg/tglimiter"
	"github.com/lodthe/bdaytracker-go/vk"
)

// General keeps the general fields for a session
type General struct {
	Bot      *telegram.Bot
	Executor *tglimiter.Executor

	StateRepo *state.Repo
	Config    conf.Config
	VKCli     *vk.Client
}
