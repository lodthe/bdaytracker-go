package application

import (
	"encoding/json"

	"github.com/petuhovskiy/telegram"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/conf"
)

func setupBot(config conf.Telegram) *telegram.Bot {
	opts := &telegram.Opts{}
	opts.Middleware = func(handler telegram.RequestHandler) telegram.RequestHandler {
		return func(methodName string, request interface{}) (json.RawMessage, error) {
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
