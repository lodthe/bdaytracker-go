package application

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lodthe/bdaytracker-go/internal/conf"
)

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
