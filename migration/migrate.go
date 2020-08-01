package migration

import (
	"github.com/jinzhu/gorm"

	"github.com/lodthe/bdaytracker-go/tg"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(
		tg.State{},
	)

	return nil
}
