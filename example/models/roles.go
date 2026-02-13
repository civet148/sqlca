package models

type Role struct {
	ID    uint
	Name  string `gorm:"uniqueIndex;size:64;"`
	Users []User `gorm:"many2many:user_roles;"`
}
