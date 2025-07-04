package models

import "time"

type SessionNote struct {
    SessionID uint
    NoteID    uint   `gorm:"primaryKey"`
    Notes     string
    Timestamp time.Time
    UserID uint
    
    User User
    Session Session
}
