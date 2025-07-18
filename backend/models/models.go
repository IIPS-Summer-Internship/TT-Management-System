package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Users []User
}

type User struct {
	gorm.Model
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	RoleID       uint
	Role         Role `gorm:"foreignKey:RoleID"`
}

type Faculty struct {
	gorm.Model
	UserID    uint   `gorm:"unique;not null"`
	User      User   `gorm:"foreignKey:UserID"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Lectures  []Lecture
}

type Course struct {
	gorm.Model
	Code     string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	Subjects []Subject
	Batches  []Batch
}

type Subject struct {
	gorm.Model
	Code     string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`
	Lectures []Lecture
}

type Batch struct {
	gorm.Model
	CourseID  uint   `gorm:"not null"`
	Course    Course `gorm:"foreignKey:CourseID"`
	EntryYear int    `gorm:"not null"`
	Sections  []Section
}

type Section struct {
	gorm.Model
	BatchID uint   `gorm:"not null"`
	Batch   Batch  `gorm:"foreignKey:BatchID"`
	Name    string `gorm:"not null"`
}

type Room struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Timetable struct {
	gorm.Model
	BatchID   uint     `gorm:"not null"`
	Batch     Batch    `gorm:"foreignKey:BatchID"`
	SectionID *uint    // Pointer for nullable foreign key
	Section   *Section `gorm:"foreignKey:SectionID"`
	CourseID  uint     `gorm:"not null"`
	Course    Course   `gorm:"foreignKey:CourseID"`
	RoomID    uint     `gorm:"not null"`
	Room      Room     `gorm:"foreignKey:RoomID"`
	Semester  int      `gorm:"not null"`
	CreatedBy uint     `gorm:"not null"`
	Creator   User     `gorm:"foreignKey:CreatedBy"`
	Lectures  []Lecture
	Sessions  []Session
}

type Timeslot struct {
	gorm.Model
	DayOfWeek int       `gorm:"not null"` // 1 for Sunday, 2 for Monday, etc.
	StartTime time.Time `gorm:"type:time"`
	EndTime   time.Time `gorm:"type:time"`
	Lectures  []Lecture
}

type Lecture struct {
	TimetableID uint   `gorm:"primaryKey"`
	TimeslotID  uint   `gorm:"primaryKey"`
	SubjectID   uint   `gorm:"not null"`
	FacultyID   uint   `gorm:"not null"`
	Room        string `gorm:"not null"`
	Timetable   Timetable
	Timeslot    Timeslot
	Subject     Subject
	Faculty     Faculty
}

type Session struct {
	gorm.Model
	TimetableID uint      `gorm:"not null"`
	TimeslotID  uint      `gorm:"not null"`
	Date        time.Time `gorm:"type:date"`
	Status      string
	Timetable   Timetable
	Timeslot    Timeslot
	Notes       []SessionNote
}

type SessionNote struct {
	gorm.Model
	SessionID uint    `gorm:"not null"`
	Session   Session `gorm:"foreignKey:SessionID"`
	EnteredBy uint    `gorm:"not null"`
	User      User    `gorm:"foreignKey:EnteredBy"`
	Notes     string  `gorm:"type:text"`
	Timestamp time.Time
}
