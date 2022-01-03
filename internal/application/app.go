package application

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/conf"
	"github.com/lodthe/bdaytracker-go/internal/mailing"
	"github.com/lodthe/bdaytracker-go/internal/tghandle"
	"github.com/lodthe/bdaytracker-go/internal/tglimiter"
	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
	"github.com/lodthe/bdaytracker-go/internal/vk"
)

type Application struct {
	ctx context.Context

	db           *sqlx.DB
	stateRepo    *tgstate.Repo
	bot          *telegram.Bot
	issuer       *usersession.Issuer
	updCollector *tghandle.UpdatesCollector

	cfg *conf.Config
}

func NewApplication(ctx context.Context) (*Application, error) {
	config := conf.Read()

	db, err := setupDatabaseConnection(config.DB)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup db conn")
	}

	err = applyMigrations(db, config.DB)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply migrations")
	}

	stateRepo := tgstate.NewRepository(db)
	bot := setupBot(config.Telegram)
	vkCli := vk.NewClient(config.VK.Token)
	telegramExecutor := tglimiter.NewExecutor()
	sessionIssuer := usersession.NewIssuer(&config, bot, telegramExecutor, vkCli, stateRepo)
	collector := tghandle.NewUpdatesCollector(sessionIssuer)

	return &Application{
		ctx:          ctx,
		cfg:          &config,
		db:           db,
		stateRepo:    stateRepo,
		bot:          bot,
		issuer:       sessionIssuer,
		updCollector: collector,
	}, nil
}

func (a *Application) Run() {
	ch, err := updates.StartPolling(a.bot, telegram.GetUpdatesRequest{
		Offset: 0,
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to start the polling")
	}

	go a.updCollector.Start(ch)

	go mailing.NewService(&a.cfg.Mailing, a.stateRepo, a.issuer).Run(a.ctx)

	logrus.Info("application started")
}

func (a *Application) Shutdown() {
	err := a.db.Close()
	if err != nil {
		logrus.WithError(err).Error("failed to close db connection")
	}

	logrus.Info("application finished")
}
