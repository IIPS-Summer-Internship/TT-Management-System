package models

type Role struct {
    RoleID uint `gorm:"primaryKey"`
    Role_Name   string
}
