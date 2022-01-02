package main

import (
	"context"
	"encoding/json"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/lodthe/bdaytracker-go/tg/state"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"

	"github.com/lodthe/bdaytracker-go/tg/callback"
	"github.com/lodthe/bdaytracker-go/tg/notifications"
	"github.com/lodthe/bdaytracker-go/tg/sessionstorage"
	"github.com/lodthe/bdaytracker-go/tg/tglimiter"
	"github.com/lodthe/bdaytracker-go/vk"

	"github.com/lodthe/bdaytracker-go/conf"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/handle"

	"github.com/jmoiron/sqlx"
)

func main() {
	setupLogging()
	config := conf.Read()

	globalContext, cancel := context.WithCancel(context.Background())

	db, err := setupDatabaseConnection(config.DB)
	if err != nil {
		logrus.WithError(err).Fatal("failed to setup db conn")
	}

	err = applyMigrations(db, config.DB)
	if err != nil {
		logrus.WithError(err).Fatal("failed to apply migrations")
	}

	stateRepo := state.NewRepository(db)

	bot := setupBot(config.Telegram)
	callback.Init()

	vkCli := vk.NewClient(config.VK.Token)

	telegramExecutor := tglimiter.NewExecutor()

	general := tg.General{
		Bot:       bot,
		Executor:  telegramExecutor,
		StateRepo: stateRepo,
		Config:    config,
		VKCli:     vkCli,
	}

	// Start getting updates from Telegram
	ch, err := updates.StartPolling(bot, telegram.GetUpdatesRequest{
		Offset: 0,
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to start the polling")
	}

	sessionStorage := sessionstorage.NewStorage()

	go notifications.NewService(stateRepo, &general, sessionStorage).Run(globalContext)

	collector := handle.NewUpdatesCollector(sessionStorage)
	collector.Start(general, ch)

	cancel()
}

func setupLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func setupDatabaseConnection(config conf.DB) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", config.PostgresDSN)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(config.MaxConnectionLifetime)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)

	return db, nil
}

func applyMigrations(db *sqlx.DB, config conf.DB) error {
	migrationDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create postgres instance")
	}

	manager, err := migrate.NewWithDatabaseInstance("file://"+config.MigrationPath, config.DabataseName, migrationDriver)
	if err != nil {
		return errors.Wrap(err, "failed to create migration manager")
	}

	err = manager.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}

func setupBot(config conf.Telegram) *telegram.Bot {
	opts := &telegram.Opts{}
	opts.Middleware = func(handler telegram.RequestHandler) telegram.RequestHandler {
		return func(methodName string, request interface{}) (json.RawMessage, error) {
			logrus.WithFields(logrus.Fields{
				"request": request,
				"method":  methodName,
			}).Debug("a telegram bot request")

			j, err := handler(methodName, request)

			if err != nil {
				logrus.WithFields(logrus.Fields{
					"request": request,
					"method":  methodName,
				}).WithError(err).Error("telegram bot request failed")
			}

			return j, err
		}
	}

	return telegram.NewBotWithOpts(config.BotToken, opts)
}
