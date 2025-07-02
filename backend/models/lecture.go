package models

type Lecture struct {
    TimetableID uint `gorm:"primaryKey"`
    TimeslotID  uint `gorm:"primaryKey"`
    SubjectID   uint
    Room_Name   string
    FacultyID   uint
    RoomID 		uint

    Timetable Timetable
    Timeslot  Timeslot
    Subject   Subject
    Faculty   Faculty
    Room Room 
}
