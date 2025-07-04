package models

type Timeslot struct {
    TimeslotID uint `gorm:"primaryKey"`
    DayOfWeek  int
    StartTime  string
    EndTime    string
}
