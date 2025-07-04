package models

type Faculty struct {
    FacultyID uint `gorm:"primaryKey"`
    UserID    uint
    Name string

    User User
}
