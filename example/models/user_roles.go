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
	UserId uint64 `gorm:"primaryKey;column:user_id;type:bigint unsigned;not null;default:0;index:fk_user_roles_user,priority:1;comment:用户ID" json:"user_id"`
	RoleId uint64 `gorm:"primaryKey;column:role_id;type:bigint unsigned;not null;default:0;index:fk_user_roles_role,priority:1;comment:角色ID" json:"role_id"`
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
