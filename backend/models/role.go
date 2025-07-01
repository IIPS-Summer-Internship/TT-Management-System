package models

type Role struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"uniqueIndex;not null"`
	Description   string
	LevelPriority int `gorm:"not null"`
}
