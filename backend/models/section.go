package models

type Section struct {
    SectionID uint `gorm:"primaryKey"`
    Name      string
    BatchID   uint

    Batch Batch
}
