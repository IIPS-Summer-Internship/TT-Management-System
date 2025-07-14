package models

type User struct {
    UserID       uint `gorm:"primaryKey"`
    Role_Name 		string	`gorm:"not null"`

    Email        string		`gorm:"uniqueIndex;not null"`
    PasswordHash string		`gorm:"not null"`
    RoleID       uint	`gorm:"not null"`
    
    Role Role `gorm:"foreignKey:RoleID;references:RoleID"`
    
}
