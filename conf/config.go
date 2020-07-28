package conf

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Telegram Telegram
	DB       DB
}

type Telegram struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type DB struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	Name     string `env:"DB_NAME" envDefault:"bdaytracker"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

func Read() Config {
	conf := Config{}

	if err := env.Parse(&conf); err != nil {
		log.WithError(err).Fatal("failed to read the config")
	}

	return conf
}
