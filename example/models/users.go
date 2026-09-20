package models

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3"
)

const TableNameUsers = "users"

const (
	UsersColumn_Id        = "id"
	UsersColumn_UserName  = "user_name"
	UsersColumn_State     = "state"
	UsersColumn_Email     = "email"
	UsersColumn_ExtraData = "extra_data"
)

const (
	USERS_COLUMN_ID         = "id"
	USERS_COLUMN_CREATED_AT = "created_at"
	USERS_COLUMN_UPDATED_AT = "updated_at"
	USERS_COLUMN_USER_NAME  = "user_name"
	USERS_COLUMN_EMAIL      = "email"
)

type User struct {
	Id            uint64                     `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	UserName      string                     `json:"user_name" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name,priority:1;null;" sqlca:"nullable"`
	State         int8                       `json:"state" db:"state" gorm:"column:state;type:tinyint(1);default:0;null;" sqlca:"nullable"`
	Email         string                     `json:"email" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email,priority:1;null;" sqlca:"nullable"`
	ExtraData     struct{}                   `json:"extra_data" db:"extra_data" gorm:"column:extra_data;type:json;null;" sqlca:"nullable"`
	Roles         []*Role                    `json:"roles,omitempty" db:"-" gorm:"many2many:user_roles;"` // 用户角色列表
	Profile       UserProfile                `json:"profile,omitempty" db:"-" gorm:"foreignKey:UserId;"`  // 用户资料明细
	BaseModel     `json:"-" gorm:"embedded"` // 基础模型(嵌入结构体)
	sqlca.LrvTree `json:"-" gorm:"embedded"` // 左右值树(嵌入结构体)
}

func (do User) DatabaseName() string { return "test" }

func (do User) TableName() string { return TableNameUsers }

func (do User) GetId() uint64 { return do.Id }

func (do User) GetUserName() string { return do.UserName }

func (do User) GetState() int8 { return do.State }

func (do User) GetEmail() string { return do.Email }

func (do User) GetExtraData() struct{} { return do.ExtraData }

func (do *User) SetId(v uint64) { do.Id = v }

func (do *User) SetUserName(v string) { do.UserName = v }

func (do *User) SetState(v int8) { do.State = v }

func (do *User) SetEmail(v string) { do.Email = v }

func (do *User) AfterQueryData(db *sqlca.Engine, ok bool) error {
	log.Infof("query ok [%v]", ok)
	do.isExist = true
	return nil
}

func (do *User) SetExtraData(v struct{}) { do.ExtraData = v }
