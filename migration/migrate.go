package migration

import (
	"github.com/jinzhu/gorm"

	"github.com/lodthe/bdaytracker-go/tg/state"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(
		state.StateDB{},
	)

	return nil
}
