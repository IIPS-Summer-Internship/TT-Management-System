package models

type Lecture struct {
    TimetableID uint `gorm:"primaryKey"`
    TimeslotID  uint `gorm:"primaryKey"`
    SubjectID   uint `gorm:"not null"`
    Room_Name   string `gorm:"not null"`
    FacultyID   uint  `gorm:"not null"`
    RoomID 		uint  `gorm:"not null"`

    Timetable Timetable  `gorm:"foreignKey:TimetableID;references:TimetableID"`
    Timeslot  Timeslot  `gorm:"foreignKey:TimeslotID;references:TimeslotID"`
    Subject   Subject   `gorm:"foreignKey:SubjectID;references:SubjectID"`
    Faculty   Faculty	`gorm:"foreignKey:FacultyID;references:FacultyID"`
    Room Room 	`gorm:"foreignKey:RoomID;references:RoomID"`
}
