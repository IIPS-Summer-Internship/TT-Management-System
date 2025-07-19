package models

type Batch struct {
	BatchID   uint   `gorm:"primaryKey"`
	EntryYear int    `gorm:"not null"`
	CourseID  uint   `gorm:"not null"`
	Course    Course `gorm:"foreignKey:CourseID;references:ID"`
}
