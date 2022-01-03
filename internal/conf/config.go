package conf

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AssetsPath string `env:"ASSETS_PATH" envDefault:"./assets"`

	Telegram Telegram
	VK       VK
	DB       DB
}

type Telegram struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type VK struct {
	Token string `env:"VK_TOKEN"`
}

type DB struct {
	PostgresDSN string `env:"DB_POSTGRES_DSN,required" envDefault:"host=localhost port=5435 user=user password=password dbname=bdaytracker sslmode=disable"`

	DabataseName  string `env:"DB_NAME,required" envDefault:"bdaytracker"`
	MigrationPath string `env:"DB_MIGRATION_PATH" envDefault:"./migrations"`

	MaxOpenConnections    int           `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"10"`
	MaxIdleConnections    int           `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	MaxConnectionLifetime time.Duration `env:"DB_MAX_CONNECTION_LIFETIME" envDefault:"5m"`
}

func Read() Config {
	conf := Config{}

	if err := env.Parse(&conf); err != nil {
		logrus.WithError(err).Fatal("failed to read the config")
	}

	return conf
}
