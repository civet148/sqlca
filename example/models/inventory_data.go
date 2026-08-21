package models

import (
	"time"

	"github.com/civet148/sqlca/v3"
)

const TableNameInventoryData = "inventory_data"

const (
	InventoryDataColumn_Id           = "id"
	InventoryDataColumn_IsFrozen     = "is_frozen"
	InventoryDataColumn_Name         = "name"
	InventoryDataColumn_SerialNo     = "serial_no"
	InventoryDataColumn_Quantity     = "quantity"
	InventoryDataColumn_Price        = "price"
	InventoryDataColumn_Location     = "location"
	InventoryDataColumn_ProductExtra = "product_extra"
	InventoryDataColumn_CreateId     = "create_id"
	InventoryDataColumn_CreateName   = "create_name"
	InventoryDataColumn_UpdateId     = "update_id"
	InventoryDataColumn_UpdateName   = "update_name"
)

type InventoryData struct {
	Id           uint64            `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	IsFrozen     FrozenState       `json:"is_frozen" db:"is_frozen" gorm:"column:is_frozen;type:tinyint(1);default:0;" sqlca:"nullable"`
	Name         string            `json:"name" db:"name" gorm:"column:name;type:varchar(255);default:null;comment:产品：名称；不能为空;" sqlca:"nullable"`                                       //产品：名称；不能为空
	SerialNo     string            `json:"serial_no" db:"serial_no" gorm:"column:serial_no;type:varchar(64);index:i_serial_no,priority:1;default:null;comment:产品序列号;" sqlca:"nullable"` //产品序列号
	Quantity     float64           `json:"quantity" db:"quantity" gorm:"column:quantity;type:decimal(16,3);default:0.000;" sqlca:"nullable"`
	Price        *float64          `json:"price" db:"price" gorm:"column:price;type:decimal(16,2);default:0.00;" sqlca:"nullable"`
	Location     sqlca.Point       `json:"location" db:"location" gorm:"column:location;type:point;default:null;" sqlca:"nullable"`
	ProductExtra *ProductExtraData `json:"product_extra" db:"product_extra" gorm:"column:product_extra;type:json;default:null;" sqlca:"nullable"`
	CreateId     uint64            `json:"create_id" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;" sqlca:"nullable"`
	CreateName   string            `json:"create_name" db:"create_name" gorm:"column:create_name;type:varchar(64);default:null;" sqlca:"nullable"`
	UpdateId     uint64            `json:"update_id" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;" sqlca:"nullable"`
	UpdateName   string            `json:"update_name" db:"update_name" gorm:"column:update_name;type:varchar(64);default:null;" sqlca:"nullable"`
	BaseModel
}

func (do InventoryData) DatabaseName() string {
	return "test"
}

func (do InventoryData) TableName() string {
	return TableNameInventoryData
}

func (do InventoryData) GetId() uint64 { return do.Id }

func (do InventoryData) GetIsFrozen() FrozenState { return do.IsFrozen }

func (do InventoryData) GetName() string { return do.Name }

func (do InventoryData) GetSerialNo() string { return do.SerialNo }

func (do InventoryData) GetQuantity() float64 { return do.Quantity }

func (do InventoryData) GetPrice() *float64 { return do.Price }

func (do InventoryData) GetLocation() sqlca.Point { return do.Location }

func (do InventoryData) GetProductExtra() *ProductExtraData { return do.ProductExtra }

func (do InventoryData) GetCreateId() uint64 { return do.CreateId }

func (do InventoryData) GetCreateName() string { return do.CreateName }

func (do InventoryData) GetUpdateId() uint64 { return do.UpdateId }

func (do InventoryData) GetUpdateName() string { return do.UpdateName }

func (do *InventoryData) SetId(v uint64) { do.Id = v }

func (do *InventoryData) SetIsFrozen(v FrozenState) { do.IsFrozen = v }

func (do *InventoryData) SetName(v string) { do.Name = v }

func (do *InventoryData) SetSerialNo(v string) { do.SerialNo = v }

func (do *InventoryData) SetQuantity(v float64) { do.Quantity = v }

func (do *InventoryData) SetPrice(v *float64) { do.Price = v }

func (do *InventoryData) SetLocation(v sqlca.Point) { do.Location = v }

func (do *InventoryData) SetProductExtra(v *ProductExtraData) { do.ProductExtra = v }

func (do *InventoryData) SetCreateId(v uint64) { do.CreateId = v }

func (do *InventoryData) SetCreateName(v string) { do.CreateName = v }

func (do *InventoryData) SetUpdateId(v uint64) { do.UpdateId = v }

func (do InventoryData) GetCreatedAt() time.Time { return do.CreatedAt }

func (do InventoryData) GetUpdatedAt() time.Time { return do.UpdatedAt }

func (do *InventoryData) SetCreatedAt(v time.Time) { do.CreatedAt = v }

func (do *InventoryData) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }

func (do *InventoryData) SetUpdateName(v string) { do.UpdateName = v }
