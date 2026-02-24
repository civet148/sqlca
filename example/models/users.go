package models

// 主模型
type User struct {
	BaseModel
	UserName string `gorm:"uniqueIndex;size:32;"`
	Email    string `gorm:"uniqueIndex;size:64;"`
	// 一对一
	Profile UserProfile `gorm:"foreignKey:UserID"`
	// 多对多
	Roles []Role `gorm:"many2many:user_roles;"`
}
