package models

import "time"

type SessionNote struct {
    SessionID uint	`gorm:"unique;not null"`
    NoteID    uint   `gorm:"primaryKey"`
    Notes     string	`gorm:"not null"`
    Timestamp time.Time	`gorm:"type:time;not null"`
    UserID uint	`gorm:"not null"`
    
    User User	`gorm:"foreignKey:UserID;references:UserID"`
    Session Session	`gorm:"foreignKey:SessionID;references:SessionID"`
}
