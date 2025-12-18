package models

const TableNameInventoryOut = "`inventory_out`" //出库主表

const (
	INVENTORY_OUT_COLUMN_ID          = "id"
	INVENTORY_OUT_COLUMN_CREATE_ID   = "create_id"
	INVENTORY_OUT_COLUMN_CREATE_NAME = "create_name"
	INVENTORY_OUT_COLUMN_CREATE_TIME = "create_time"
	INVENTORY_OUT_COLUMN_UPDATE_ID   = "update_id"
	INVENTORY_OUT_COLUMN_UPDATE_NAME = "update_name"
	INVENTORY_OUT_COLUMN_UPDATE_TIME = "update_time"
	INVENTORY_OUT_COLUMN_IS_DELETED  = "is_deleted"
	INVENTORY_OUT_COLUMN_DELETE_TIME = "delete_time"
	INVENTORY_OUT_COLUMN_PRODUCT_ID  = "product_id"
	INVENTORY_OUT_COLUMN_ORDER_NO    = "order_no"
	INVENTORY_OUT_COLUMN_USER_ID     = "user_id"
	INVENTORY_OUT_COLUMN_USER_NAME   = "user_name"
	INVENTORY_OUT_COLUMN_QUANTITY    = "quantity"
	INVENTORY_OUT_COLUMN_WEIGHT      = "weight"
	INVENTORY_OUT_COLUMN_REMARK      = "remark"
)

type InventoryOut struct {
	Id         uint64  `gorm:"column:id;index;"`          //主键ID
	CreateId   uint64  `gorm:"column:create_id;index;"`   //创建人ID
	CreateName string  `gorm:"column:create_name;index;"` //创建人姓名
	CreateTime string  `gorm:"column:create_time;index;"` //创建时间
	UpdateId   uint64  `gorm:"column:update_id;index;"`   //更新人ID
	UpdateName string  `gorm:"column:update_name;index;"` //更新人姓名
	UpdateTime string  `gorm:"column:update_time;index;"` //更新时间
	IsDeleted  int8    `gorm:"column:is_deleted;index;"`  //删除状态(0: 未删除 1: 已删除)
	DeleteTime string  `gorm:"column:delete_time;index;"` //删除时间
	ProductId  uint64  `gorm:"column:product_id;index;"`  //产品ID
	OrderNo    string  `gorm:"column:order_no;index;"`    //出库单号
	UserId     uint64  `gorm:"column:user_id;index;"`     //收货人ID
	UserName   string  `gorm:"column:user_name;index;"`   //收货人姓名
	Quantity   float64 `gorm:"column:quantity;index;"`    //数量
	Weight     float64 `gorm:"column:weight;index;"`      //净重
	Remark     string  `gorm:"column:remark;index;"`      //备注
}

func (do *InventoryOut) GetId() uint64          { return do.Id }
func (do *InventoryOut) SetId(v uint64)         { do.Id = v }
func (do *InventoryOut) GetCreateId() uint64    { return do.CreateId }
func (do *InventoryOut) SetCreateId(v uint64)   { do.CreateId = v }
func (do *InventoryOut) GetCreateName() string  { return do.CreateName }
func (do *InventoryOut) SetCreateName(v string) { do.CreateName = v }
func (do *InventoryOut) GetCreateTime() string  { return do.CreateTime }
func (do *InventoryOut) SetCreateTime(v string) { do.CreateTime = v }
func (do *InventoryOut) GetUpdateId() uint64    { return do.UpdateId }
func (do *InventoryOut) SetUpdateId(v uint64)   { do.UpdateId = v }
func (do *InventoryOut) GetUpdateName() string  { return do.UpdateName }
func (do *InventoryOut) SetUpdateName(v string) { do.UpdateName = v }
func (do *InventoryOut) GetUpdateTime() string  { return do.UpdateTime }
func (do *InventoryOut) SetUpdateTime(v string) { do.UpdateTime = v }
func (do *InventoryOut) GetIsDeleted() int8     { return do.IsDeleted }
func (do *InventoryOut) SetIsDeleted(v int8)    { do.IsDeleted = v }
func (do *InventoryOut) GetDeleteTime() string  { return do.DeleteTime }
func (do *InventoryOut) SetDeleteTime(v string) { do.DeleteTime = v }
func (do *InventoryOut) GetProductId() uint64   { return do.ProductId }
func (do *InventoryOut) SetProductId(v uint64)  { do.ProductId = v }
func (do *InventoryOut) GetOrderNo() string     { return do.OrderNo }
func (do *InventoryOut) SetOrderNo(v string)    { do.OrderNo = v }
func (do *InventoryOut) GetUserId() uint64      { return do.UserId }
func (do *InventoryOut) SetUserId(v uint64)     { do.UserId = v }
func (do *InventoryOut) GetUserName() string    { return do.UserName }
func (do *InventoryOut) SetUserName(v string)   { do.UserName = v }
func (do *InventoryOut) GetQuantity() float64   { return do.Quantity }
func (do *InventoryOut) SetQuantity(v float64)  { do.Quantity = v }
func (do *InventoryOut) GetWeight() float64     { return do.Weight }
func (do *InventoryOut) SetWeight(v float64)    { do.Weight = v }
func (do *InventoryOut) GetRemark() string      { return do.Remark }
func (do *InventoryOut) SetRemark(v string)     { do.Remark = v }
