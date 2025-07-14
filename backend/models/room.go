package models

type Room struct {
    RoomID uint `gorm:"primaryKey"`
    Room_Name   string  `gorm:"uniqueIndex;not null"`
}
