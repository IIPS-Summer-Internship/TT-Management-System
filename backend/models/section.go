package models

type Section struct {
    SectionID uint `gorm:"primaryKey"`
    Name      string `gorm:"not null"`
    BatchID   uint	`gorm:"not null"`

    Batch Batch	`gorm:"foreignKey:BatchID;references:BatchID"`
}
