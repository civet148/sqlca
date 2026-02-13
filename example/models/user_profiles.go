package models

type UserProfile struct {
	ID      uint   // 唯一ID
	UserID  uint   `gorm:"uniqueIndex"` // 确保一对一
	Avatar  string `gorm:"size:512"`    // 头像URL
	Address string `gorm:"size:128"`    // 家庭住址
}
