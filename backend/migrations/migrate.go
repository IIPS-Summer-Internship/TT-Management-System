package migrations

import (
	"tms-server/config"
	"tms-server/models"
)

func Migrate() error {
	err := config.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Faculty{},
		&models.Course{},
		&models.Batch{},
		&models.Section{},
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
