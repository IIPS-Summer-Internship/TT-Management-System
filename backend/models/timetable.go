package models

type Timetable struct {
    TimetableID uint `gorm:"primaryKey"`
    BatchID     uint
    RoomID      uint
    Semester    uint
    CourseID    uint
    UserID   uint
    
    User User
    Batch   Batch
    Room    Room
    Course  Course
}
