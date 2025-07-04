package models

type Batch struct {
    BatchID    uint `gorm:"primaryKey"`
    EntryYear  int
    CourseID   uint

    Course Course
}
