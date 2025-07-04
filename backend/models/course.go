package models

type Course struct {
    CourseID uint `gorm:"primaryKey"`
    Code     string 
    Name     string 
}
