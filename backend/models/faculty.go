package models

type Faculty struct {
    FacultyID uint `gorm:"primaryKey"`
    UserID    uint `gorm:"unique;not null"`
    Name      string `gorm:"not null"`

    User User `gorm:"foreignKey:UserID;references:UserID"`
}

