package models

import (
	"time"
)

const TableNameUserRoles = "user_roles"

const (
	UserRolesColumn_UserId = "user_id"
	UserRolesColumn_RoleId = "role_id"
)

type UserRole struct {
	UserId uint64 `json:"user_id" db:"user_id" gorm:"column:user_id;type:bigint unsigned;;default:0;"`
	RoleId uint64 `json:"role_id" db:"role_id" gorm:"column:role_id;type:bigint unsigned;index:fk_user_roles_role,priority:1;;default:0;"`
	BaseModel
}

func (do UserRole) DatabaseName() string {
	return "test"
}

func (do UserRole) TableName() string {
	return TableNameUserRoles
}

func (do UserRole) GetUserId() uint64 { return do.UserId }

func (do UserRole) GetRoleId() uint64 { return do.RoleId }

func (do *UserRole) SetUserId(v uint64) { do.UserId = v }

func (do UserRole) GetCreatedAt() time.Time { return do.CreatedAt }

func (do UserRole) GetUpdatedAt() time.Time { return do.UpdatedAt }

func (do *UserRole) SetCreatedAt(v time.Time) { do.CreatedAt = v }

func (do *UserRole) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }

func (do *UserRole) SetRoleId(v uint64) { do.RoleId = v }
