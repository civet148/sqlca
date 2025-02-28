# 简介
sqlca 是一个基于Go语言的ORM框架，它提供了一种简单的方式来生成数据库表模型，并支持多种数据库类型，如MySQL、PostgreSQL、Opengauss、MS-SQLServer、Sqlite v3等。
内置雪花算法生成主键ID、SSH隧道连接以及防SQL注入功能。支持各种数据库聚合方法和联表查询，例如: Sum/Max/Avg/Min/Count/GroupBy/Having/OrderBy/Limit等等。
同时将常用的操作符进行了包装，例如等于Eq、大于Gt、小于Lt等等，简化操作代码。其中And和Or方法既支持常规的字符串格式化（含占位符?方式），同时也支持map类型传参作为WHERE/AND/OR条件进行查询和以及更新。

## sqlca与gorm差异

- sqlca不支持通过数据模型自动生成创建/更新时间（可由MySQL等数据库创建表时设置为由数据库自动维护生成/更新时间），当数据库自动维护创建/更新时间时，可通过 `sqlca:"readonly"` 标签将数据字段设置为只读

- sqlca由Model方法调用后，会自动克隆一个对象，后续所有的操作均不影响宿主对象。每当一个完整语句执行完毕（例如调用Query/Update/Delete方法后)，db对象会清理掉所有的查询/更新条件。

- 对于model结构嵌套的差异

```go
type CommonFields struct {
	Id         int64  `db:"id"`          //数据库主键ID
	UpdateTime string `db:"update_time"` //更新时间
	CreateTime string `db:"create_time"` //创建时间
}

type User struct {
    CommonFields CommonFields //没有db标签
    Name string `db:"name"` //姓名
    Gender int32 `db:"gender"` //性别
}
```
对于上面的User结构，对CommonFields的处理sqlca插入和查询跟gorm保持一致，都是把id/update_time/create_time字段作为跟name，gender平级的字段处理。

```go
type ExtraData struct {
	Address     string `json:"address"` //家庭住址
	Email       string `json:"email"`   //电子邮箱地址
}

type User struct {
    Id          int64       `db:"id"`          //数据库主键ID
    UpdateTime  string      `db:"update_time"` //更新时间
    CreateTime  string      `db:"create_time"` //创建时间
    Name        string      `db:"name"`        //姓名
    Gender      int32       `db:"gender"`      //性别
    ExtraData   ExtraData   `db:"extra_data"`  //额外数据
}
```
对于上面的User结构，ExtraData成员变量因为有db标签，sqlca把ExtraData作为user表的一个字段进行处理，插入时把ExtraData序列化为JSON文本存入extra_data字段。查询时反序列化到ExtraData结构中。
而gorm把ExtraData作为外键处理。

## sqlca标签说明

- `sqlca:"readonly"` 只读标签，指定该标签的字段插入和更新操作均不参与
- `sqlca:"isnull"`  允许为空标签，指定该标签的字段允许为空(数据库字段允许为NULL)

## db2go工具
[db2go](https://github.com/civet148/db2go) 是一个支持从MySQL、PostgreSQL、Opengauss数据库导出表结构到.go文件或.proto文件的命令行工具。支持将表字段指定为自定义类型并生成model文件和dao文件。

# 快速开始

## 支持数据库类型

- **MySQL**

```text
"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4"
"mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4"
```

- **PostgreSQL**
```text
"postgres://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public"
```
- **Opengauss**
```text
"opengauss://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public"
```
- **MS-SQLServer**
```text
"mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false"
```
- **Sqlite v3**
```text
"sqlite:///var/lib/test.db"
```

## 数据库表模型生成

- 创建数据库

```sql
CREATE DATABASE `test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `test`;


CREATE TABLE `inventory_data` (
                                  `id` bigint unsigned NOT NULL COMMENT '产品ID',
                                  `create_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
                                  `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '创建人姓名',
                                  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `update_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新人ID',
                                  `update_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '更新人姓名',
                                  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                  `is_frozen` tinyint(1) NOT NULL DEFAULT '0' COMMENT '冻结状态(0: 未冻结 1: 已冻结)',
                                  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '产品名称',
                                  `serial_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '产品编号',
                                  `quantity` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '产品库存',
                                  `price` decimal(16,2) NOT NULL DEFAULT '0.00' COMMENT '产品均价',
                                  `product_extra` json DEFAULT NULL COMMENT '产品附带数据(JSON文本)',
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存数据表';

CREATE TABLE `inventory_in` (
                                `id` bigint unsigned NOT NULL COMMENT '主键ID',
                                `create_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
                                `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '创建人姓名',
                                `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `update_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新人ID',
                                `update_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '更新人姓名',
                                `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除状态(0: 未删除 1: 已删除)',
                                `delete_time` datetime DEFAULT NULL COMMENT '删除时间',
                                `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
                                `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '入库单号',
                                `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '交货人ID',
                                `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交货人姓名',
                                `quantity` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '数量',
                                `weight` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '净重',
                                `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
                                PRIMARY KEY (`id`) USING BTREE,
                                UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='入库主表';

CREATE TABLE `inventory_out` (
                                 `id` bigint unsigned NOT NULL COMMENT '主键ID',
                                 `create_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
                                 `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '创建人姓名',
                                 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `update_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新人ID',
                                 `update_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '更新人姓名',
                                 `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                 `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除状态(0: 未删除 1: 已删除)',
                                 `delete_time` datetime DEFAULT NULL COMMENT '删除时间',
                                 `product_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '产品ID',
                                 `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '出库单号',
                                 `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '收货人ID',
                                 `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '收货人姓名',
                                 `quantity` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '数量',
                                 `weight` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '净重',
                                 `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
                                 PRIMARY KEY (`id`) USING BTREE,
                                 UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='出库主表';

```

- **安装db2go工具**

```shell
$ go install github.com/civet148/db2go@latest
```

- **自动生成go代码脚本**

```bat
@echo off

rem 设置只读字段
set READ_ONLY="create_time, update_time"
rem 数据模型(models)和数据库操作对象(dao)文件输出基础目录
set OUT_DIR=.
rem 数据模型包名(数据模型文件目录名)
set PACK_NAME="models"
rem 指定某表的某字段为指定类型,多个表字段以英文逗号分隔（例如：user.create_time=time.Time表示指定user表create_time字段为time.Time类型; 如果不指定表名则所有表的create_time字段均为time.Time类型；支持第三方包类型，例如：user.weight=github.com/shopspring/decimal.Decimal）
set SPEC_TYPES="inventory_data.product_extra=ProductExtraData"
rem 指定其他orm的标签和值(以空格分隔)
set COMMON_TAGS="id=gorm:\"primarykey\" create_time=gorm:\"autoCreateTime\" update_time=gorm:\"autoUpdateTime\""
set DEPLOY_SQL="test.sql"

rem 判断本地系统是否已安装db2go工具，没有则进行安装
echo "searching db2go.exe ..."
echo "--------------------------------------------"
where db2go.exe
echo "--------------------------------------------"

IF "%errorlevel%" == "0" (
    echo db2go already installed.
) ELSE (
    echo db2go not found in system %%PATH%%, installing...
    go install github.com/civet148/db2go@latest
    If "%errorlevel%" == "0" (
        echo db2go install successfully.
    ) ELSE (
        rem 安装失败: Linux/Mac请安装gcc工具链，Windows系统可以安装msys64进行源码编译或通过链接直接下载二进制(最新版本v2.13 https://github.com/civet148/release/tree/master/db2go/v2)
        echo ERROR: Linux/Mac please install gcc tool-chain and windows download from https://github.com/civet148/release/tree/master/db2go/v2 (latest version is v2.13)
    )
)

rem ---------------------- 导出数据库表结构-------------------------
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
db2go   --url %DSN_URL% --out %OUT_DIR% --spec-type %SPEC_TYPES% --package %PACK_NAME%  --common-tags %COMMON_TAGS% --readonly %READ_ONLY% --export %DEPLOY_SQL%
gofmt -w %OUT_DIR%/%PACK_NAME%

pause
```

- **生成的代码示例**

```go
// Code generated by db2go. DO NOT EDIT.
// https://github.com/civet148/db2go

package models

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
	Id           uint64           `json:"id" db:"id" gorm:"primarykey"`                                         //产品ID
	CreateId     uint64           `json:"create_id" db:"create_id" `                                            //创建人ID
	CreateName   string           `json:"create_name" db:"create_name" `                                        //创建人姓名
	CreateTime   string           `json:"create_time" db:"create_time" gorm:"autoCreateTime" sqlca:"readonly"`  //创建时间
	UpdateId     uint64           `json:"update_id" db:"update_id" `                                            //更新人ID
	UpdateName   string           `json:"update_name" db:"update_name" `                                        //更新人姓名
	UpdateTime   string           `json:"update_time" db:"update_time" gorm:"autoUpdateTime" sqlca:"readonly"`  //更新时间
	IsFrozen     int8             `json:"is_frozen" db:"is_frozen" `                                            //冻结状态(0: 未冻结 1: 已冻结)
	Name         string           `json:"name" db:"name" `                                                      //产品名称
	SerialNo     string           `json:"serial_no" db:"serial_no" `                                            //产品编号
	Quantity     float64          `json:"quantity" db:"quantity" `                                              //产品库存
	Price        float64          `json:"price" db:"price" `                                                    //产品均价
	ProductExtra ProductExtraData `json:"product_extra" db:"product_extra" sqlca:"isnull"`                      //产品附带数据(JSON文本)
}

func (do *InventoryData) GetId() uint64                      { return do.Id }
func (do *InventoryData) SetId(v uint64)                     { do.Id = v }
func (do *InventoryData) GetCreateId() uint64                { return do.CreateId }
func (do *InventoryData) SetCreateId(v uint64)               { do.CreateId = v }
func (do *InventoryData) GetCreateName() string              { return do.CreateName }
func (do *InventoryData) SetCreateName(v string)             { do.CreateName = v }
func (do *InventoryData) GetCreateTime() string              { return do.CreateTime }
func (do *InventoryData) SetCreateTime(v string)             { do.CreateTime = v }
func (do *InventoryData) GetUpdateId() uint64                { return do.UpdateId }
func (do *InventoryData) SetUpdateId(v uint64)               { do.UpdateId = v }
func (do *InventoryData) GetUpdateName() string              { return do.UpdateName }
func (do *InventoryData) SetUpdateName(v string)             { do.UpdateName = v }
func (do *InventoryData) GetUpdateTime() string              { return do.UpdateTime }
func (do *InventoryData) SetUpdateTime(v string)             { do.UpdateTime = v }
func (do *InventoryData) GetIsFrozen() int8                  { return do.IsFrozen }
func (do *InventoryData) SetIsFrozen(v int8)                 { do.IsFrozen = v }
func (do *InventoryData) GetName() string                    { return do.Name }
func (do *InventoryData) SetName(v string)                   { do.Name = v }
func (do *InventoryData) GetSerialNo() string                { return do.SerialNo }
func (do *InventoryData) SetSerialNo(v string)               { do.SerialNo = v }
func (do *InventoryData) GetQuantity() float64               { return do.Quantity }
func (do *InventoryData) SetQuantity(v float64)              { do.Quantity = v }
func (do *InventoryData) GetPrice() float64                  { return do.Price }
func (do *InventoryData) SetPrice(v float64)                 { do.Price = v }
func (do *InventoryData) GetProductExtra() ProductExtraData  { return do.ProductExtra }
func (do *InventoryData) SetProductExtra(v ProductExtraData) { do.ProductExtra = v }


```

## 连接数据库

```golang
package main

import (
    "github.com/civet148/log"
    "github.com/civet148/sqlca/v2"
)

const (
	//MysslDSN = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4"
    MysqlDSN = "mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4"
	//PostgresDSN  = "postgres://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public")
	//GaussDSN  = "opengauss://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public")
	//MssqlDSN  = "mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false")
	//SqliteDSN  = "sqlite:///var/lib/test.db")
)

func main() {
	var err error
	var db *sqlca.Engine
	var options = &sqlca.Options{
		Debug: true, //是否开启调试日志输出
		Max:   150,  //最大连接数
		Idle:  5,    //空闲连接数
		//DefaultLimit: 100,  //默认查询条数限制
		SnowFlake: &sqlca.SnowFlake{ //雪花算法配置（不使用可以赋值nil）
			NodeId: 1, //雪花算法节点ID 1-1023
		},
	}
	db, err = sqlca.NewEngine(MysqlDSN, options)
	if err != nil {
		log.Errorf("connect database error: %s", err)
		return
	}
	_ = db
}

```

# 数据库CURD示例

## 单条插入

```go
func InsertSingle(db *sqlca.Engine) error {
	
	now := time.Now().Format("2006-01-02 15:04:05")
	var do = &models.InventoryData{
		Id:         uint64(db.NewID()),
		CreateId:   1,
		CreateName: "admin",
		CreateTime: now,
		UpdateId:   1,
		UpdateName: "admin",
		UpdateTime: now,
		IsFrozen:   0,
		Name:       "齿轮",
		SerialNo:   "SNO_001",
		Quantity:   1000,
		Price:      10.5,
	}

	var err error
	/*
		INSERT INTO inventory_data (`id`,`create_id`,`create_name`,`create_time`,`update_id`,`update_name`,`update_time`,`is_frozen`,`name`,`serial_no`,`quantity`,`price`,`product_extra`)
		VALUES ('1859078192380252161','1','admin','2024-11-20 11:35:55','1','admin','2024-11-20 11:35:55','0','轮胎','SNO_002','2000','210','{}')
	*/
	_, err = db.Model(&do).Insert()
	if err != nil {
		return log.Errorf("数据插入错误: %s", err)
	}
	return nil
}
```

## 批量插入

```go
func InsertBatch(db *sqlca.Engine) error {
now := time.Now().Format("2006-01-02 15:04:05")
var dos = []*models.InventoryData{
        {
            Id:         uint64(db.NewID()),
            CreateId:   1,
            CreateName: "admin",
            CreateTime: now,
            UpdateId:   1,
            UpdateName: "admin",
            UpdateTime: now,
            IsFrozen:   0,
            Name:       "齿轮",
            SerialNo:   "SNO_001",
            Quantity:   1000,
            Price:      10.5,
            ProductExtra: models.ProductExtraData{
                SpecsValue: "齿数：32",
                AvgPrice:   sqlca.NewDecimal(30.8),
            },
        },
        {
            Id:         uint64(db.NewID()),
            CreateId:   1,
            CreateName: "admin",
            CreateTime: now,
            UpdateId:   1,
            UpdateName: "admin",
            UpdateTime: now,
            IsFrozen:   0,
            Name:       "轮胎",
            SerialNo:   "SNO_002",
            Quantity:   2000,
            Price:      210,
            ProductExtra: models.ProductExtraData{
                SpecsValue: "17英寸",
                AvgPrice:   sqlca.NewDecimal(450.5),
            },
        },
    }

    var err error
    /*
        INSERT INTO inventory_data
            (`id`,`create_id`,`create_name`,`create_time`,`update_id`,`update_name`,`update_time`,`is_frozen`,`name`,`serial_no`,`quantity`,`price`,`product_extra`)
        VALUES
            ('1867379968636358656','1','admin','2024-12-13 09:24:13','1','admin','2024-12-13 09:24:13','0','齿轮','SNO_001','1000','10.5','{\"avg_price\":\".8\",\"specs_value\":\"齿数：32\"}'),
            ('1867379968636358657','1','admin','2024-12-13 09:24:13','1','admin','2024-12-13 09:24:13','0','轮胎','SNO_002','2000','210','{\"avg_price\":\"450.5\",\"specs_value\":\"17英寸\"}')
    */
    _, err = db.Model(&dos).Insert()
    if err != nil {
        return log.Errorf("数据插入错误: %s", err)
    }
    return nil
}
```

## 普通查询带LIMIT限制
```go
func QueryLimit(db *sqlca.Engine) error {
	
    var err error
    var count int64
    var dos []*models.InventoryData
    
    //SELECT * FROM inventory_data ORDER BY create_time DESC LIMIT 1000
    count, err = db.Model(&dos).
        Limit(1000).
        Desc("create_time").
        Query()
    if err != nil {
        return log.Errorf("数据查询错误：%s", err)
    }
    log.Infof("查询结果数据条数: %d", count)
    return nil
}

```

## 查询无数据则报错

```go
func QueryErrNotFound(db *sqlca.Engine) error {
	
	var err error
	var count int64
	var do *models.InventoryData //如果取数对象是切片则不报错

	//SELECT * FROM inventory_data WHERE id=1899078192380252160
	count, err = db.Model(&do).Id(1899078192380252160).Find()
	if err != nil {
		if errors.Is(err, sqlca.ErrRecordNotFound) {
			return log.Errorf("根据ID查询数据库记录无结果：%s", err)
		}
		return log.Errorf("数据库错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}
```

## 分页查询

```go
func QueryByPage(db *sqlca.Engine) error {
	
	var err error
	var count, total int64
	var dos []*models.InventoryData

	//SELECT  * FROM inventory_data WHERE 1=1 ORDER BY create_time DESC LIMIT 0,20
	count, total, err = db.Model(&dos).
		Page(1, 20).
		Desc("create_time").
		QueryEx()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果条数: %d 数据库总数：%v", count, total)
	return nil
}

```

## 多条件查询

```go
func QueryByCondition(db *sqlca.Engine) error {
	
	var err error
	var count int64
	var dos []*models.InventoryData
	//SELECT * FROM inventory_data WHERE `quantity` > 0 and is_frozen=0 AND create_time >= '2024-10-01 11:35:14' ORDER BY create_time DESC
	count, err = db.Model(&dos).
		Gt("quantity", 0).
		Eq("is_frozen", 0).
		Gte("create_time", "2024-10-01 11:35:14").
		Desc("create_time").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}
```

## 查询JSON内容字段到数据对象
```go
/*
   models.InventoryData对象的ProductExtra是一个跟数据库JSON内容对应的结构体, 数据库中product_extra字段是json类型或text文本

    type ProductExtraData struct {
        AvgPrice   sqlca.Decimal `json:"avg_price"`   //均价
        SpecsValue string        `json:"specs_value"` //规格
    }
*/
func QueryWithJsonColumn(db *sqlca.Engine) error {
    var err error
    var do models.InventoryData
    var id = uint64(1867379968636358657)
    
    /*
        SELECT * FROM inventory_data WHERE id=1867379968636358657
    
        +-----------------------+-----------------------+-----------------------+------------------------------------------------+
        | id	                | name	| serial_no	    | quantity	| price	    |                 product_extra                  |
        +-----------------------+-------+---------------+-----------+-----------+------------------------------------------------+
        | 1867379968636358657	| 轮胎  	| SNO_002		| 2000.000 	| 210.00	| {"avg_price": "450.5", "specs_value": "17英寸"} |
        +------------------------------------------------------------------------------------------------------------------------+
    */
    _, err = db.Model(&do).
                Table("inventory_data").
                Select("id", "name", "serial_no", "quantity","price", "product_extra").
                Id(id).
                Find()
    if err != nil {
        return log.Errorf("数据查询错误：%s", err)
    }
    log.Infof("ID: %v 数据：%+v", id, do)
    /*
        2024-12-18 15:15:03.560732 PID:64764 [INFO] {goroutine 1} <main.go:373 QueryWithJsonColumn()> ID: 1867379968636358657 数据：{Id:1867379968636358657 Name:轮胎 SerialNo:SNO_002 Quantity:2000 Price:210 ProductExtra:{AvgPrice:450.5 SpecsValue:17英寸}}
    */
    return nil
}
```

## 原生SQL查询

```go
func QueryRawSQL(db *sqlca.Engine) error {
	
	var rows []*models.InventoryData
	var sb = sqlca.NewStringBuilder()

	//SELECT * FROM inventory_data  WHERE is_frozen =  '0' AND quantity > '10'

	sb.Append("SELECT * FROM %s", "inventory_data")
	sb.Append("WHERE is_frozen = ?", 0)
	sb.Append("AND quantity > ?", 10)
	strQuery := sb.String()
	_, err := db.Model(&rows).QueryRaw(strQuery)
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	return nil
}
```

## 原生SQL执行

```go
func ExecRawSQL(db *sqlca.Engine) error {
    var sb = sqlca.NewStringBuilder()
    
    //UPDATE inventory_data SET quantity = '10' WHERE id=1867379968636358657
    sb.Append("UPDATE inventory_data")
    sb.Append("SET quantity = ?", 10)
    sb.Append("WHERE id = ?", 1867379968636358657)
    
    strQuery := sb.String()
    affectedRows, lastInsertId, err := db.Model(nil).ExecRaw(strQuery)
    if err != nil {
        return log.Errorf("数据查询错误：%s", err)
    }
    log.Infof("受影响的行数：%d 最后插入的ID：%d", affectedRows, lastInsertId)
    return nil
}

```

## 查询带多个OR条件(map类型)

```go
func QueryOr(db *sqlca.Engine) error {
	
	var err error
	var count int64
	var dos []*models.InventoryData

	//SELECT * FROM inventory_data WHERE create_id=1 AND name = '配件' OR serial_no = 'SNO_001' ORDER BY create_time DESC
	count, err = db.Model(&dos).
		And("create_id = ?", 1).
		Or("name = ?", "配件").
		Or("serial_no = ?", "SNO_001").
		Desc("create_time").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)

	//SELECT * FROM inventory_data WHERE create_id=1 AND is_frozen = 0 AND quantity > 0 AND (name = '配件' OR serial_no = 'SNO_001') ORDER BY create_time DESC
	var andConditions = make(map[string]interface{})
	var orConditions = make(map[string]interface{})

	andConditions["create_id"] = 1      //create_id = 1
	andConditions["is_frozen"] = 0      //is_frozen = 0
    andConditions["quantity > ?"] = 0   //quantity > 0

	orConditions["name = ?"] = "配件"
	orConditions["serial_no = ?"] = "SNO_001"

	count, err = db.Model(&dos).
		And(andConditions).
		Or(orConditions).
		Desc("create_time").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}
```

## 分组查询

```go
func QueryByGroup(db *sqlca.Engine) error {
	
	var err error
	var count int64
	var dos []*models.InventoryData
	/*
		SELECT  create_id, SUM(quantity) AS quantity
		FROM inventory_data
		WHERE 1=1 AND quantity>'0' AND is_frozen='0' AND create_time>='2024-10-01 11:35:14'
		GROUP BY create_id
	*/
	count, err = db.Model(&dos).
		Select("create_id", "SUM(quantity) AS quantity").
		Gt("quantity", 0).
		Eq("is_frozen", 0).
		Gte("create_time", "2024-10-01 11:35:14").
		GroupBy("create_id").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}
```


## 联表查询

```go
type Product struct {
    ProductId       uint64  `db:"product_id"`
    ProductName     string  `db:"product_name"`
    Quantity        float64 `db:"quantity"`
    Weight          float64 `db:"weight"` 
}
func QueryJoins(db *sqlca.Engine) error {
	
	/*
		SELECT a.id as product_id, a.name AS product_name, b.quantity, b.weight
		FROM inventory_data a
		LEFT JOIN inventory_in b
		ON a.id=b.product_id
		WHERE a.quantity > 0 AND a.is_frozen=0 AND a.create_time>='2024-10-01 11:35:14'
	*/
	var dos []*Product
	count, err := db.Model(&dos).
		Select("a.id as product_id", "a.name AS product_name", "b.quantity", "b.weight").
		Table("inventory_data a").
		LeftJoin("inventory_in b").
		On("a.id=b.product_id").
		Gt("a.quantity", 0).
		Eq("a.is_frozen", 0).
		Gte("a.create_time", "2024-10-01 11:35:14").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}
```
## 普通变量取值查询

```go
func QueryByNormalVars(db *sqlca.Engine) error {
	
    var err error
    var name, serialNo string
    var id = uint64(1859078192380252160)
    //SELECT name, serial_no FROM inventory_data WHERE id=1859078192380252160
    _, err = db.Model(&name, &serialNo).
                Table("inventory_data").
                Select("name, serial_no").
                Id(id).
                Find()
    if err != nil {
        return log.Errorf("数据查询错误：%s", err)
    }
    log.Infof("数据ID: %v name=%s serial_no=%s", id, name, serialNo)
	
    var ids []uint64
    //SELECT id FROM inventory_data LIMIT 10
    _, err = db.Model(&ids).
                Table("inventory_data").
                Select("id").
                Limit(10).
                Query()
    if err != nil {
        return log.Errorf("数据查询错误：%s", err)
    }
    return nil
}
```

## 数据更新

- **通过数据模型对象更新数据**

```go
/*
[数据更新]

SELECT * FROM inventory_data  WHERE `id`='1858759254329004032'
UPDATE inventory_data SET `quantity`='2300' WHERE `id`='1858759254329004032'
*/
func UpdateByModel(db *sqlca.Engine) error {
	
    var err error
    var do *models.InventoryData
    var id = uint64(1858759254329004032)
    _, err = db.Model(&do).Id(id).Find() //Find方法如果是单条记录没找到则提示ErrNotFound错误（Query方法不会报错）
    if err != nil {
        return log.Errorf("数据查询错误：%s", err)
    }
    
    do.Quantity = 2300 //更改库存
    _, err = db.Model(do).Select("quantity").Update()
    if err != nil {
        return log.Errorf("更新错误：%s", err)
    }
    return nil
}
```

- **通过变量/常量更新数据**

```go
/*
[通过普通变量更新数据]
*/
func UpdateByVars(db *sqlca.Engine) error {
	
    var err error
    var id = uint64(1858759254329004032)
    var quantity = 2300 //更改库存数
	
	//UPDATE inventory_data SET `quantity`='2300' WHERE `id`='1858759254329004032'
    _, err = db.Model(&quantity).Table("inventory_data").Id(id).Select("quantity").Update()
    if err != nil {
        return log.Errorf("更新错误：%s", err)
    }
    //UPDATE inventory_data SET `quantity`='2300' WHERE `id`='1858759254329004032'
    _, err = db.Model(2300).Table("inventory_data").Id(id).Select("quantity").Update()
    if err != nil {
        return log.Errorf("更新错误：%s", err)
    }
    return nil
}
```

- **通过map进行数据更新**

```go
func UpdateByMap(db *sqlca.Engine) error {
    var err error
    var updates = map[string]interface{}{
        "quantity": 2100, //更改库存
        "Price":    300,  //更改价格
    }
    //UPDATE inventory_data SET `quantity`='2100',`price`=300 WHERE `id`='1858759254329004032'
    _, err = db.Model(&updates).Table("inventory_data").Id(1858759254329004032).Update()
    if err != nil {
        return log.Errorf("更新错误：%s", err)
    }
    return nil
}
```

- **删除操作**

```go
func DeleteById(db *sqlca.Engine) error {
    var err error
    var id = uint64(1859078192380252160)
    //DELETE inventory_data WHERE `id`='1859078192380252160'
    _, err = db.Model(models.InventoryData{}).Id(id).Delete()
    if err != nil {
        return log.Errorf("更新错误：%s", err)
    }
    log.Infof("删除ID%v数据成功", id)
    return nil
}
```


## 事务处理

```go
func Transaction(db *sqlca.Engine) error {

	/*
		-- TRANSACTION BEGIN

			INSERT INTO inventory_in (`user_id`,`quantity`,`remark`,`create_id`,`user_name`,`weight`,`create_time`,`update_name`,`is_deleted`,`product_id`,`id`,`create_name`,`update_id`,`update_time`,`order_no`) VALUES ('3','20','产品入库','1','lazy','200.3','2024-11-27 11:35:14','admin','0','1858759254329004032','1861614736295071744','admin','1','2024-11-27 1114','202407090000001')
			SELECT * FROM inventory_data  WHERE `id`='1858759254329004032'
			UPDATE inventory_data SET `quantity`='2320' WHERE `id`='1858759254329004032'

		-- TRANSACTION END
	*/

	now := time.Now().Format("2006-01-02 15:04:05")
	tx, err := db.TxBegin()
	if err != nil {
		return log.Errorf("开启事务失败：%s", err)
	}
	defer tx.TxRollback()

	productId := uint64(1858759254329004032)
	strOrderNo := time.Now().Format("20060102150405.000000000")
	//***************** 执行事务操作 *****************
	quantity := float64(20)
	weight := float64(200.3)
	_, err = tx.Model(&models.InventoryIn{
		Id:         uint64(db.NewID()),
		CreateId:   1,
		CreateName: "admin",
		CreateTime: now,
		UpdateId:   1,
		UpdateName: "admin",
		UpdateTime: now,
		ProductId:  productId,
		OrderNo:    strOrderNo,
		UserId:     3,
		UserName:   "lazy",
		Quantity:   quantity,
		Weight:     weight,
		Remark:     "产品入库",
	}).Insert()
	if err != nil {
		return log.Errorf("数据插入错误: %s", err)
	}
	var inventoryData = &models.InventoryData{}
	_, err = tx.Model(&inventoryData).Id(productId).Find() //Find方法如果是单条记录没找到则提示ErrNotFound错误（Query方法不会报错）
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	inventoryData.Quantity += quantity
	_, err = tx.Model(&inventoryData).Id(productId).Select("quantity").Update()
	if err != nil {
		return log.Errorf("更新错误：%s", err)
	}
	//***************** 提交事务 *****************
	err = tx.TxCommit()
	if err != nil {
		return log.Errorf("提交事务失败：%s", err)
	}
	return nil
}

```

## 事务处理封装

```go
func TransactionWrapper(db *sqlca.Engine) error {

    /*
	   -- TRANSACTION BEGIN

	   	INSERT INTO inventory_in (`user_id`,`quantity`,`remark`,`create_id`,`user_name`,`weight`,`create_time`,`update_name`,`is_deleted`,`product_id`,`id`,`create_name`,`update_id`,`update_time`,`order_no`) VALUES ('3','20','产品入库','1','lazy','200.3','2024-11-27 11:35:14','admin','0','1858759254329004032','1861614736295071744','admin','1','2024-11-27 1114','202407090000002')
	   	SELECT * FROM inventory_data  WHERE `id`='1858759254329004032'
	   	UPDATE inventory_data SET `quantity`='2320' WHERE `id`='1858759254329004032'

	   -- TRANSACTION END
	*/
	strOrderNo := time.Now().Format("20060102150405.000000000")
	err := db.TxFunc(func(tx *sqlca.Engine) error {
		var err error
		productId := uint64(1858759254329004032)
		now := time.Now().Format("2006-01-02 15:04:05")

		//***************** 执行事务操作 *****************
		quantity := float64(20)
		weight := float64(200.3)
		_, err = tx.Model(&models.InventoryIn{
			Id:         uint64(db.NewID()),
			CreateId:   1,
			CreateName: "admin",
			CreateTime: now,
			UpdateId:   1,
			UpdateName: "admin",
			UpdateTime: now,
			ProductId:  productId,
			OrderNo:    strOrderNo,
			UserId:     3,
			UserName:   "lazy",
			Quantity:   quantity,
			Weight:     weight,
			Remark:     "产品入库",
		}).Insert()
		if err != nil {
			return log.Errorf("数据插入错误: %s", err)
		}
		var inventoryData = &models.InventoryData{}
		_, err = tx.Model(&inventoryData).Id(productId).Find() //Find方法如果是单条记录没找到则提示ErrNotFound错误（Query方法不会报错）
		if err != nil {
			return log.Errorf("数据查询错误：%s", err)
		}
		inventoryData.Quantity += quantity
		_, err = tx.Model(&inventoryData).Id(productId).Select("quantity").Update()
		if err != nil {
			return log.Errorf("更新错误：%s", err)
		}
		return nil
	})

	//***************** 事务处理结果 *****************
	if err != nil {
		return log.Errorf("事务失败：%s", err)
	}
	return nil
}
```
## 其他方法说明

### Table
设置数据库表名，通过Model方法传参时默认将结构体名称的小写蛇形命名作为表名，当传入的结构体名称跟实际表名不一致时需要明确用Table方法指定表名

### NearBy

```go
//数据库表restaurant对应模型结构
type Restaurant struct {
    Id          uint64  `db:"id"`       //主键ID
    Lng         float64 `db:"lng"`      //经度
    Lat         float64 `db:"lat"`      //纬度
    Name        string  `db:"name"`     //餐馆名称
}

//附近的餐馆和距离结构定义
type RestaurantLocation struct {
	Id          uint64  `db:"id"`       //主键ID
	Lng         float64 `db:"lng"`      //经度
	Lat         float64 `db:"lat"`      //纬度
	Name        string  `db:"name"`     //餐馆名称
	Distance    float64 `db:"distance"` //距离（米）
}

func QueryNearBy(db *sqlca.Engine) error {
    var dos []*RestaurantLocation
    //查询指定坐标点，查询距离小于1000米内的餐馆（查询出的距离取名distance）
    _, err := db.Model(&dos).Table("restaurant").NearBy("lng", "lat", "distance", 114.0545429, 22.5445741, 1000).Query()
    if err != nil {
        return logx.Error(err.Error())
    }
    return nil
}

```
### GeoHash

给定坐标点，查询GEO HASH

### Like

```go
func QueryLike(db *sqlca.Engine) error {
    //SELECT * FROM inventory_data WHERE `serial_no` LIKE '%0001%'
    _, err := db.Model(&models.InventoryData{}).LIKE(serial_no, "0001").Find()
    if err != nil {
        return logx.Error(err.Error())
    }
	return nil
}
```

### SlowQuery 
开启或关闭慢查询日志，默认关闭，开启后，会记录超过规定时间（毫秒ms）的sql语句，并输出到控制台。

### QueryJson
将查询结果转换为json字符串，并返回。

### NewID
当调用NewEngine时，指定SnowFlake选项后，可以用NewID生成一个雪花ID

### NewFromTx
传入一个tx对象，并返回一个Engine对象，用于在事务中执行sql操作。

### ForUpdate
在查询语句中添加FOR UPDATE关键字，用于查询时锁定记录，避免并发修改。仅用于MySQL数据库。

### LockShareMode
在查询语句中添加 LOCK IN SHARE MODE关键字。仅用于MySQL数据库。


### JSON查询方法

#### **JsonExpr**
MySQL数据库构造JSON查询表达式，用于查询JSON字段。

#### **JsonEqual**
MySQL数据库构造JSON等于查询表达式，用于查询JSON字段。

#### **JsonGreater**
MySQL数据库构造JSON大于查询表达式，用于查询JSON字段。

#### **JsonLess**
MySQL数据库构造JSON小于查询表达式，用于查询JSON字段。

#### **JsonGreaterEqual**
MySQL数据库构造JSON大于等于查询表达式，用于查询JSON字段。

#### **JsonLessEqual**
MySQL数据库构造JSON小于等于查询表达式，用于查询JSON字段。

#### **JsonContainArray**
MySQL数据库构造JSON包含数组查询表达式，用于查询JSON字段。
