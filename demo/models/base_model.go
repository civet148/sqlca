package models

import (
	"time"

	"github.com/civet148/sqlca/v3"
)

type BaseModel struct {
	Id         uint64 `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                           //产品ID
	CreateTime string `json:"create_time,omitempty" db:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP;autoCreateTime"` //创建时间
	CreateId   uint64 `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id"`                                                //创建人ID
	CreateName string `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name"`                                          //创建人姓名
	UpdateId   uint64 `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id"`                                                //更新人ID
	UpdateName string `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name"`                                          //更新人姓名
	UpdateTime string `json:"update_time,omitempty" db:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP;autoUpdateTime"` //更新时间
	isExist    bool   `db:"-"`                                                                                                           //数据是否在数据库中存在
}

// 是否为数据库中存在的数据
func (b BaseModel) IsExist() bool {
	return b.isExist
}

// 是否为新的数据
func (b BaseModel) IsNew() bool {
	return !b.isExist
}

func NowTime() string {
	return time.Now().Format(time.DateTime)
}

func (do *BaseModel) BeforeQueryData(db *sqlca.Engine) error {
	return nil
}

func (do *BaseModel) AfterQueryData(db *sqlca.Engine) error {
	if do.Id != 0 {
		do.isExist = true
	}
	return nil
}

func (do *BaseModel) BeforeCreateData(db *sqlca.Engine) error {
	do.CreateTime = NowTime()
	do.UpdateTime = NowTime()
	return nil
}

func (do *BaseModel) AfterCreateData(db *sqlca.Engine) error {
	return nil
}

func (do *BaseModel) BeforeUpdateData(db *sqlca.Engine) error {
	return nil
}

func (do *BaseModel) AfterUpdateData(db *sqlca.Engine) error {
	return nil
}

func (do *BaseModel) BeforeDeleteData(db *sqlca.Engine) error {
	return nil
}

func (do *InventoryData) AfterDeleteData(db *sqlca.Engine) error {
	return nil
}
