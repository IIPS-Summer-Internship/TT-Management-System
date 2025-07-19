package models

type Subject struct {
	SubjectID uint   `gorm:"primaryKey"`
	CourseID  uint   `gorm:"not null"`
	Code      string `gorm:"not null"`
	Name      string `gorm:"not null"`

	Course Course `gorm:"foreignKey:CourseID;references:CourseID"`
}
