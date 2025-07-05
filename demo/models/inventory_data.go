package models

type FrozenState int

const (
	FrozenState_False = 0
	FrozenState_Ture  = 1
)

func (s FrozenState) String() string {
	switch s {
	case FrozenState_Ture:
		return "True"
	case FrozenState_False:
		return "False"
	}
	return "<FrozenState_Unknown>"
}

const TableNameInventoryData = "`inventory_data`" //库存数据表

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
)

type InventoryData struct {
	Id           uint64            `json:"id" db:"id" gorm:"primarykey"`                                        //产品ID
	CreateId     uint64            `json:"create_id" db:"create_id" `                                           //创建人ID
	CreateName   string            `json:"create_name" db:"create_name" `                                       //创建人姓名
	CreateTime   string            `json:"create_time" db:"create_time" gorm:"autoCreateTime" sqlca:"readonly"` //创建时间
	UpdateId     uint64            `json:"update_id" db:"update_id" `                                           //更新人ID
	UpdateName   string            `json:"update_name" db:"update_name" `                                       //更新人姓名
	UpdateTime   string            `json:"update_time" db:"update_time" gorm:"autoUpdateTime" sqlca:"readonly"` //更新时间
	IsFrozen     FrozenState       `json:"is_frozen" db:"is_frozen" `                                           //冻结状态(0: 未冻结 1: 已冻结)
	Name         string            `json:"name" db:"name" `                                                     //产品名称
	SerialNo     string            `json:"serial_no" db:"serial_no" `                                           //产品编号
	Quantity     float64           `json:"quantity" db:"quantity" `                                             //产品库存
	Price        *float64          `json:"price" db:"price" `                                                   //产品均价
	ProductExtra *ProductExtraData `json:"product_extra" db:"product_extra" sqlca:"isnull"`                     //产品附带数据(JSON文本)
	Nullable     string            `json:"nullable" db:"nullable" `                                             //可为空字段
}

func (do *InventoryData) GetId() uint64                       { return do.Id }
func (do *InventoryData) SetId(v uint64)                      { do.Id = v }
func (do *InventoryData) GetCreateId() uint64                 { return do.CreateId }
func (do *InventoryData) SetCreateId(v uint64)                { do.CreateId = v }
func (do *InventoryData) GetCreateName() string               { return do.CreateName }
func (do *InventoryData) SetCreateName(v string)              { do.CreateName = v }
func (do *InventoryData) GetCreateTime() string               { return do.CreateTime }
func (do *InventoryData) SetCreateTime(v string)              { do.CreateTime = v }
func (do *InventoryData) GetUpdateId() uint64                 { return do.UpdateId }
func (do *InventoryData) SetUpdateId(v uint64)                { do.UpdateId = v }
func (do *InventoryData) GetUpdateName() string               { return do.UpdateName }
func (do *InventoryData) SetUpdateName(v string)              { do.UpdateName = v }
func (do *InventoryData) GetUpdateTime() string               { return do.UpdateTime }
func (do *InventoryData) SetUpdateTime(v string)              { do.UpdateTime = v }
func (do *InventoryData) GetIsFrozen() FrozenState            { return do.IsFrozen }
func (do *InventoryData) SetIsFrozen(v FrozenState)           { do.IsFrozen = v }
func (do *InventoryData) GetName() string                     { return do.Name }
func (do *InventoryData) SetName(v string)                    { do.Name = v }
func (do *InventoryData) GetSerialNo() string                 { return do.SerialNo }
func (do *InventoryData) SetSerialNo(v string)                { do.SerialNo = v }
func (do *InventoryData) GetQuantity() float64                { return do.Quantity }
func (do *InventoryData) SetQuantity(v float64)               { do.Quantity = v }
func (do *InventoryData) GetPrice() *float64                  { return do.Price }
func (do *InventoryData) SetPrice(v *float64)                 { do.Price = v }
func (do *InventoryData) GetProductExtra() *ProductExtraData  { return do.ProductExtra }
func (do *InventoryData) SetProductExtra(v *ProductExtraData) { do.ProductExtra = v }

func (do *InventoryData) TableName() string {
	return "inventory_data"
}
