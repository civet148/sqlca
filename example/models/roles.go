package models

const TableNameRoles = "roles" //

const (
	ROLES_COLUMN_ID   = "id"
	ROLES_COLUMN_NAME = "name"
)

type Role struct {
	BaseModel
	Name string `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(64);uniqueIndex:idx_roles_name;" sqlca:"isnull"` //
}

func (do Role) TableName() string { return "roles" }

func (do Role) GetId() uint64   { return do.Id }
func (do Role) GetName() string { return do.Name }

func (do *Role) SetId(v uint64)   { do.Id = v }
func (do *Role) SetName(v string) { do.Name = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
