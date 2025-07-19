package models

type Course struct {
	CourseID uint   `gorm:"primaryKey"`
	Code     string `gorm:"not null"`
	Name     string `gorm:"not null"`
}
