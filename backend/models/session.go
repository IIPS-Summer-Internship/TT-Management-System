package models

import "time"

type Session struct {
    SessionID   uint `gorm:"primaryKey"`
    TimetableID uint
    TimeslotID  uint
    Date        time.Time
    Status      string

    Timetable Timetable
    Timeslot  Timeslot
}
