package models

import "github.com/civet148/sqlca/v3"

const TableNameInventoryData = "inventory_data" //产品库存数据表

const (
	INVENTORY_DATA_COLUMN_ID            = "id"
	INVENTORY_DATA_COLUMN_CREATE_ID     = "create_id"
	INVENTORY_DATA_COLUMN_CREATE_NAME   = "create_name"
	INVENTORY_DATA_COLUMN_CREATE_TIME   = "create_time"
	INVENTORY_DATA_COLUMN_UPDATE_ID     = "update_id"
	INVENTORY_DATA_COLUMN_UPDATE_NAME   = "update_name"
	INVENTORY_DATA_COLUMN_UPDATE_TIME   = "update_time"
	INVENTORY_DATA_COLUMN_IS_FROZEN     = "is_frozen"
	INVENTORY_DATA_COLUMN_NAME          = "name"
	INVENTORY_DATA_COLUMN_SERIAL_NO     = "serial_no"
	INVENTORY_DATA_COLUMN_QUANTITY      = "quantity"
	INVENTORY_DATA_COLUMN_PRICE         = "price"
	INVENTORY_DATA_COLUMN_PRODUCT_EXTRA = "product_extra"
	INVENTORY_DATA_COLUMN_LOCATION      = "location"
)

type InventoryData struct {
	BaseModel
	Id           uint64            `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                               //产品ID
	IsFrozen     FrozenState       `json:"is_frozen,omitempty" db:"is_frozen" gorm:"column:is_frozen;type:tinyint(1);default:0;" sqlca:"isnull"`          //
	Name         string            `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(255);" sqlca:"isnull"`                                 //
	SerialNo     string            `json:"serial_no,omitempty" db:"serial_no" gorm:"column:serial_no;type:varchar(64);index:i_serial_no;" sqlca:"isnull"` //
	Quantity     float64           `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,3);default:0.000;" sqlca:"isnull"`      //
	Price        *float64          `json:"price,omitempty" db:"price" gorm:"column:price;type:decimal(16,2);default:0.00;" sqlca:"isnull"`                //
	ProductExtra *ProductExtraData `json:"product_extra,omitempty" db:"product_extra" gorm:"column:product_extra;type:text;" sqlca:"isnull"`              //
	Location     sqlca.Point       `json:"location,omitempty" db:"location" gorm:"column:location;type:point;" sqlca:"isnull"`                            //
}

func (do InventoryData) TableName() string { return "inventory_data" }

func (do InventoryData) GetId() uint64                      { return do.Id }
func (do InventoryData) GetCreateId() uint64                { return do.CreateId }
func (do InventoryData) GetCreateName() string              { return do.CreateName }
func (do InventoryData) GetCreateTime() string              { return do.CreateTime }
func (do InventoryData) GetUpdateId() uint64                { return do.UpdateId }
func (do InventoryData) GetUpdateName() string              { return do.UpdateName }
func (do InventoryData) GetUpdateTime() string              { return do.UpdateTime }
func (do InventoryData) GetIsFrozen() FrozenState           { return do.IsFrozen }
func (do InventoryData) GetName() string                    { return do.Name }
func (do InventoryData) GetSerialNo() string                { return do.SerialNo }
func (do InventoryData) GetQuantity() float64               { return do.Quantity }
func (do InventoryData) GetPrice() *float64                 { return do.Price }
func (do InventoryData) GetProductExtra() *ProductExtraData { return do.ProductExtra }
func (do InventoryData) GetLocation() sqlca.Point           { return do.Location }

func (do *InventoryData) SetId(v uint64)                      { do.Id = v }
func (do *InventoryData) SetCreateId(v uint64)                { do.CreateId = v }
func (do *InventoryData) SetCreateName(v string)              { do.CreateName = v }
func (do *InventoryData) SetCreateTime(v string)              { do.CreateTime = v }
func (do *InventoryData) SetUpdateId(v uint64)                { do.UpdateId = v }
func (do *InventoryData) SetUpdateName(v string)              { do.UpdateName = v }
func (do *InventoryData) SetUpdateTime(v string)              { do.UpdateTime = v }
func (do *InventoryData) SetIsFrozen(v FrozenState)           { do.IsFrozen = v }
func (do *InventoryData) SetName(v string)                    { do.Name = v }
func (do *InventoryData) SetSerialNo(v string)                { do.SerialNo = v }
func (do *InventoryData) SetQuantity(v float64)               { do.Quantity = v }
func (do *InventoryData) SetPrice(v *float64)                 { do.Price = v }
func (do *InventoryData) SetProductExtra(v *ProductExtraData) { do.ProductExtra = v }
func (do *InventoryData) SetLocation(v sqlca.Point)           { do.Location = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
