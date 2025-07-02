package models

type Subject struct {
    SubjectID uint `gorm:"primaryKey"`
    CourseID  uint
    Code      string
    Name      string

    Course Course
}
