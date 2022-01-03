package usersession

import (
	conf2 "github.com/lodthe/bdaytracker-go/internal/conf"
	vk2 "github.com/lodthe/bdaytracker-go/internal/vk"
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/internal/tgstate"

	"github.com/lodthe/bdaytracker-go/internal/tglimiter"
)

// General keeps the general fields for a session
type General struct {
	Bot      *telegram.Bot
	Executor *tglimiter.Executor

	StateRepo *tgstate.Repo
	Config    conf2.Config
	VKCli     *vk2.Client
}
