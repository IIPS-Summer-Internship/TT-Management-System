package models

type Timeslot struct {
    TimeslotID uint `gorm:"primaryKey"`
    DayOfWeek  int	`gorm:"not null"`  //0 - 6 days
    StartTime  string	`gorm:"not null"`	
    EndTime    string	`gorm:"not null"`
}
