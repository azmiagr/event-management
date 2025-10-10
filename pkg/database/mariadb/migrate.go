package mariadb

import (
	"event-management/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.Otp{},
		&entity.Event{},
		&entity.Registration{},
		&entity.Session{},
		&entity.Payment{},
	)
	if err != nil {
		return err
	}

	return nil
}
