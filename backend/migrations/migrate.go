package migrations

import (
	"tms-server/config"
	"tms-server/models"
)

// INFO: for UP and DOWN migration: github.com/golang-migrate/migrate/v4
func Migrate() error {
	err := config.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Faculty{},
		&models.Course{},
		&models.Batch{},
		&models.Subject{},
		&models.Room{},
		&models.Timetable{},
		&models.Timeslot{},
		&models.Lecture{},
		&models.Session{},
		&models.SessionNote{},
	)
	return err
}
