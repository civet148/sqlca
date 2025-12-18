package models

const TableNameInventoryIn = "`inventory_in`" //入库主表

const (
	INVENTORY_IN_COLUMN_ID          = "id"
	INVENTORY_IN_COLUMN_CREATE_ID   = "create_id"
	INVENTORY_IN_COLUMN_CREATE_NAME = "create_name"
	INVENTORY_IN_COLUMN_CREATE_TIME = "create_time"
	INVENTORY_IN_COLUMN_UPDATE_ID   = "update_id"
	INVENTORY_IN_COLUMN_UPDATE_NAME = "update_name"
	INVENTORY_IN_COLUMN_UPDATE_TIME = "update_time"
	INVENTORY_IN_COLUMN_IS_DELETED  = "is_deleted"
	INVENTORY_IN_COLUMN_DELETE_TIME = "delete_time"
	INVENTORY_IN_COLUMN_PRODUCT_ID  = "product_id"
	INVENTORY_IN_COLUMN_ORDER_NO    = "order_no"
	INVENTORY_IN_COLUMN_USER_ID     = "user_id"
	INVENTORY_IN_COLUMN_USER_NAME   = "user_name"
	INVENTORY_IN_COLUMN_QUANTITY    = "quantity"
	INVENTORY_IN_COLUMN_WEIGHT      = "weight"
	INVENTORY_IN_COLUMN_REMARK      = "remark"
)

type InventoryIn struct {
	Id         uint64  `gorm:"column:id"`                           //主键ID
	CreateId   uint64  `gorm:"column:create_id"`                    //创建人ID
	CreateName string  `gorm:"column:create_name"`                  //创建人姓名
	CreateTime string  `gorm:"column:create_time" sqlca:"readonly"` //创建时间
	UpdateId   uint64  `gorm:"column:update_id"`                    //更新人ID
	UpdateName string  `gorm:"column:update_name"`                  //更新人姓名
	UpdateTime string  `gorm:"column:update_time" qlca:"readonly"`  //更新时间
	IsDeleted  int8    `gorm:"column:is_deleted"`                   //删除状态(0: 未删除 1: 已删除)
	DeleteTime string  `gorm:"column:delete_time" sqlca:"isnull"`   //删除时间
	ProductId  uint64  `gorm:"column:product_id"`                   //产品ID
	OrderNo    string  `gorm:"column:order_no"`                     //入库单号
	UserId     uint64  `gorm:"column:user_id"`                      //交货人ID
	UserName   string  `gorm:"column:user_name"`                    //交货人姓名
	Quantity   float64 `gorm:"column:quantity"`                     //数量
	Weight     float64 `gorm:"column:weight"`                       //净重
	Remark     string  `gorm:"column:remark"`                       //备注
}

func (do *InventoryIn) GetId() uint64          { return do.Id }
func (do *InventoryIn) SetId(v uint64)         { do.Id = v }
func (do *InventoryIn) GetCreateId() uint64    { return do.CreateId }
func (do *InventoryIn) SetCreateId(v uint64)   { do.CreateId = v }
func (do *InventoryIn) GetCreateName() string  { return do.CreateName }
func (do *InventoryIn) SetCreateName(v string) { do.CreateName = v }
func (do *InventoryIn) GetCreateTime() string  { return do.CreateTime }
func (do *InventoryIn) SetCreateTime(v string) { do.CreateTime = v }
func (do *InventoryIn) GetUpdateId() uint64    { return do.UpdateId }
func (do *InventoryIn) SetUpdateId(v uint64)   { do.UpdateId = v }
func (do *InventoryIn) GetUpdateName() string  { return do.UpdateName }
func (do *InventoryIn) SetUpdateName(v string) { do.UpdateName = v }
func (do *InventoryIn) GetUpdateTime() string  { return do.UpdateTime }
func (do *InventoryIn) SetUpdateTime(v string) { do.UpdateTime = v }
func (do *InventoryIn) GetIsDeleted() int8     { return do.IsDeleted }
func (do *InventoryIn) SetIsDeleted(v int8)    { do.IsDeleted = v }
func (do *InventoryIn) GetDeleteTime() string  { return do.DeleteTime }
func (do *InventoryIn) SetDeleteTime(v string) { do.DeleteTime = v }
func (do *InventoryIn) GetProductId() uint64   { return do.ProductId }
func (do *InventoryIn) SetProductId(v uint64)  { do.ProductId = v }
func (do *InventoryIn) GetOrderNo() string     { return do.OrderNo }
func (do *InventoryIn) SetOrderNo(v string)    { do.OrderNo = v }
func (do *InventoryIn) GetUserId() uint64      { return do.UserId }
func (do *InventoryIn) SetUserId(v uint64)     { do.UserId = v }
func (do *InventoryIn) GetUserName() string    { return do.UserName }
func (do *InventoryIn) SetUserName(v string)   { do.UserName = v }
func (do *InventoryIn) GetQuantity() float64   { return do.Quantity }
func (do *InventoryIn) SetQuantity(v float64)  { do.Quantity = v }
func (do *InventoryIn) GetWeight() float64     { return do.Weight }
func (do *InventoryIn) SetWeight(v float64)    { do.Weight = v }
func (do *InventoryIn) GetRemark() string      { return do.Remark }
func (do *InventoryIn) SetRemark(v string)     { do.Remark = v }
