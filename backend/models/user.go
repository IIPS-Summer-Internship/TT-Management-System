package models

type User struct {
    UserID       uint `gorm:"primaryKey"`
    Role_Name 		string

    Email        string
    PasswordHash string
    RoleID       uint
    
    Role Role
    
}
