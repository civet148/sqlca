package models

type BaseModel struct {
	CreateTime string `json:"create_time,omitempty" db:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP;autoCreateTime"` //创建时间
	CreateId   uint64 `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id"`                                                //创建人ID
	CreateName string `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name"`                                          //创建人姓名
	UpdateId   uint64 `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id"`                                                //更新人ID
	UpdateName string `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name"`                                          //更新人姓名
	UpdateTime string `json:"update_time,omitempty" db:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP;autoUpdateTime"` //更新时间
}
