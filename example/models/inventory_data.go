package models

import "time"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryData = "inventory_data" //

const (
	INVENTORY_DATA_COLUMN_ID            = "id"
	INVENTORY_DATA_COLUMN_CREATED_AT    = "created_at"
	INVENTORY_DATA_COLUMN_UPDATED_AT    = "updated_at"
	INVENTORY_DATA_COLUMN_IS_FROZEN     = "is_frozen"
	INVENTORY_DATA_COLUMN_NAME          = "name"
	INVENTORY_DATA_COLUMN_SERIAL_NO     = "serial_no"
	INVENTORY_DATA_COLUMN_QUANTITY      = "quantity"
	INVENTORY_DATA_COLUMN_PRICE         = "price"
	INVENTORY_DATA_COLUMN_LOCATION      = "location"
	INVENTORY_DATA_COLUMN_PRODUCT_EXTRA = "product_extra"
)

type InventoryData struct {
	BaseModel
	CreatedAt    time.Time         `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime;index:idx_inventory_data_created_at;default:CURRENT_TIMESTAMP;"` //
	UpdatedAt    time.Time         `json:"updated_at,omitempty" db:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime;index:idx_inventory_data_updated_at;default:CURRENT_TIMESTAMP;"` //
	IsFrozen     FrozenState       `json:"is_frozen,omitempty" db:"is_frozen" gorm:"column:is_frozen;type:tinyint(1);default:0;" sqlca:"isnull"`                                                       //
	Name         string            `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(255);comment:产品：名称；不能为空;" sqlca:"isnull"`                                                           //产品：名称；不能为空
	SerialNo     string            `json:"serial_no,omitempty" db:"serial_no" gorm:"column:serial_no;type:varchar(64);index:i_serial_no;comment:产品序列号;" sqlca:"isnull"`                                //产品序列号
	Quantity     float64           `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,3);default:0.000;" sqlca:"isnull"`                                                   //
	Price        *float64          `json:"price,omitempty" db:"price" gorm:"column:price;type:decimal(16,2);default:0.00;" sqlca:"isnull"`                                                             //
	Location     sqlca.Point       `json:"location,omitempty" db:"location" gorm:"column:location;type:point;" sqlca:"isnull"`                                                                         //
	ProductExtra *ProductExtraData `json:"product_extra,omitempty" db:"product_extra" gorm:"column:product_extra;type:json;" sqlca:"isnull"`                                                           //
}

func (do InventoryData) TableName() string { return "inventory_data" }

func (do InventoryData) GetId() uint64                      { return do.Id }
func (do InventoryData) GetCreatedAt() time.Time            { return do.CreatedAt }
func (do InventoryData) GetUpdatedAt() time.Time            { return do.UpdatedAt }
func (do InventoryData) GetIsFrozen() FrozenState           { return do.IsFrozen }
func (do InventoryData) GetName() string                    { return do.Name }
func (do InventoryData) GetSerialNo() string                { return do.SerialNo }
func (do InventoryData) GetQuantity() float64               { return do.Quantity }
func (do InventoryData) GetPrice() *float64                 { return do.Price }
func (do InventoryData) GetLocation() sqlca.Point           { return do.Location }
func (do InventoryData) GetProductExtra() *ProductExtraData { return do.ProductExtra }

func (do *InventoryData) SetId(v uint64)                      { do.Id = v }
func (do *InventoryData) SetCreatedAt(v time.Time)            { do.CreatedAt = v }
func (do *InventoryData) SetUpdatedAt(v time.Time)            { do.UpdatedAt = v }
func (do *InventoryData) SetIsFrozen(v FrozenState)           { do.IsFrozen = v }
func (do *InventoryData) SetName(v string)                    { do.Name = v }
func (do *InventoryData) SetSerialNo(v string)                { do.SerialNo = v }
func (do *InventoryData) SetQuantity(v float64)               { do.Quantity = v }
func (do *InventoryData) SetPrice(v *float64)                 { do.Price = v }
func (do *InventoryData) SetLocation(v sqlca.Point)           { do.Location = v }
func (do *InventoryData) SetProductExtra(v *ProductExtraData) { do.ProductExtra = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
