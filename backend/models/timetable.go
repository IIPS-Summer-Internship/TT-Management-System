package models

type Timetable struct {
    TimetableID uint `gorm:"primaryKey"`
    BatchID     uint	`gorm:"not null"`
    RoomID      uint	`gorm:"not null"`
    Semester    uint	`gorm:"not null"`
    CourseID    uint	`gorm:"not null"`
    UserID   uint	`gorm:"not null"`
    
    User User	`gorm:"foreignKey:UserID;references:UserID"`
    Batch   Batch	`gorm:"foreignKey:BatchID;references:BatchID"`
    Room    Room	`gorm:"foreignKey:RoomID;references:RoomID"`
    Course  Course	`gorm:"foreignKey:CourseID;references:CourseID"`
}
