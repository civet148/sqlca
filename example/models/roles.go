package models

import (
	"time"
)

const TableNameRoles = "roles"

const (
	RolesColumn_Id   = "id"
	RolesColumn_Name = "name"
)

type Role struct {
	Id   uint64 `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	Name string `json:"name" db:"name" gorm:"column:name;type:varchar(64);uniqueIndex:idx_roles_name,priority:1;default:null;" sqlca:"nullable"`
	BaseModel
}

func (do Role) DatabaseName() string {
	return "test"
}

func (do Role) TableName() string {
	return TableNameRoles
}

func (do Role) GetId() uint64 { return do.Id }

func (do Role) GetName() string { return do.Name }

func (do *Role) SetId(v uint64) { do.Id = v }

func (do Role) GetCreatedAt() time.Time { return do.CreatedAt }

func (do Role) GetUpdatedAt() time.Time { return do.UpdatedAt }

func (do *Role) SetCreatedAt(v time.Time) { do.CreatedAt = v }

func (do *Role) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }

func (do *Role) SetName(v string) { do.Name = v }
