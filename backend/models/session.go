package models

import "time"

type Session struct {
	SessionID   uint      `gorm:"primaryKey"`
	TimetableID uint      `gorm:"not null"`
	TimeslotID  uint      `gorm:"not null"`
	Date        time.Time `gorm:"type:date;not null"`
	Status      string    `gorm:"not null"`

	Timetable Timetable `gorm:"foreignKey:TimetableID;references:TimetableID"`
	Timeslot  Timeslot  `gorm:"foreignKey:TimeslotID;references:TimeslotID"`
}
