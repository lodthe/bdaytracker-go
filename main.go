package main

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/migration"
	"github.com/lodthe/bdaytracker-go/tg/callback"

	"github.com/lodthe/bdaytracker-go/conf"
	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/handle"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	setupLogging()
	config := conf.Read()

	db := setupGORM(config.DB)

	bot := setupBot(config.Telegram)
	callback.Init()

	general := tg.General{
		Bot:    bot,
		DB:     db,
		Config: config,
	}

	// Start getting updates from Telegram
	ch, err := updates.StartPolling(bot, telegram.GetUpdatesRequest{
		Offset: 0,
	})
	if err != nil {
		log.WithError(err).Fatal("failed to start the polling")
	}

	collector := handle.NewUpdatesCollector()
	collector.Start(general, ch)
}

func setupLogging() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func setupGORM(config conf.DB) *gorm.DB {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s user=%s password=%s", config.Host, config.Port, config.SSLMode, config.Name, config.User, config.Password))
	if err != nil {
		log.WithError(err).Fatal("failed to open the db")
	}

	err = migration.Migrate(db)
	if err != nil {
		log.WithError(err).Fatal("failed to make migrations")
	}

	return db
}

func setupBot(config conf.Telegram) *telegram.Bot {
	opts := &telegram.Opts{}
	opts.Middleware = func(handler telegram.RequestHandler) telegram.RequestHandler {
		return func(methodName string, request interface{}) (json.RawMessage, error) {
			log.WithFields(log.Fields{
				"request": request,
				"method":  methodName,
			}).Debug("a telegram bot request")

			j, err := handler(methodName, request)

			if err != nil {
				log.WithFields(log.Fields{
					"request": request,
					"method":  methodName,
				}).WithError(err).Error("telegram bot request failed")
			}

			return j, err
		}
	}

	return telegram.NewBotWithOpts(config.BotToken, opts)
}
