package models

import (
	"time"

	"github.com/civet148/sqlca/v3"
)

type BaseModel struct {
	Id        uint64    `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                     //产品ID
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at;type:timestamp;not null;index;default:CURRENT_TIMESTAMP;autoCreatedAt"` //创建时间
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" gorm:"column:updated_at;type:timestamp;not null;index;default:CURRENT_TIMESTAMP;autoUpdatedAt"` //更新时间
	isExist   bool      `gorm:"-" db:"-"`                                                                                                                            //数据是否在数据库中存在
}

// 是否为数据库中存在的数据
func (b BaseModel) IsExist() bool {
	return b.isExist
}

// 是否为新的数据
func (b BaseModel) IsNew() bool {
	return !b.isExist
}

func NowTime() time.Time {
	return time.Now()
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
	do.CreatedAt = NowTime()
	do.UpdatedAt = NowTime()
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
