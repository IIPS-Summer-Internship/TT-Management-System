package migrations

import (
	"log"
	"tms-server/config"
	"tms-server/models"
)

// INFO: for UP and DOWN migration: github.com/golang-migrate/migrate/v4
func Migrate() error {
	err := config.DB.AutoMigrate(
		&models.Role{}, // role table for new architecture
		&models.User{},
		&models.Faculty{},
		&models.Course{},
		&models.Batch{},
		&models.Subject{},
		&models.Room{},
		&models.Lecture{},
		&models.Session{},
	)
	if err != nil {
		return err
	}
	//only migrate existing users if needed
	if config.DB.Migrator().HasColumn(&models.User{}, "role") {
		log.Println("Old role column found, but skipping migration - will be handled manually in Supabase")
	}

	return nil
}
