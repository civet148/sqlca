package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3"
	"github.com/civet148/sqlca/v3/demo/models"
	"time"
)

const (
	productId = 1939612151790440448
)

func main() {
	var err error
	var db *sqlca.Engine
	var opts = []sqlca.Option{
		sqlca.WithDebug(),
		sqlca.WithMaxConn(100),
		sqlca.WithIdleConn(5),
		//SSH tunnel config
		//sqlca.WithSSH(&sqlca.SSH{
		//	User:     "root",
		//	Password: "123456",
		//	Host:     "192.168.2.19:22",
		//}),
		//redis distribution lock config
		sqlca.WithRedisConfig(&sqlca.RedisConfig{
			Address: "192.168.1.20:6379",
		}),
		sqlca.WithSnowFlake(&sqlca.SnowFlake{
			NodeId: 1,
		}),
	}
	db, err = sqlca.NewEngine("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4", opts...)
	if err != nil {
		log.Errorf("connect database error: %s", err)
		return
	}

	requireNoError(InsertSingle(db))
	requireNoError(InsertBatch(db))
	requireNoError(UpsertSingle(db))
	requireNoError(QueryLimit(db))
	requireError(QueryErrNotFound(db))
	requireNoError(QueryByPage(db))
	requireNoError(QueryByCondition(db))
	requireNoError(QueryByGroup(db))
	requireNoError(QueryCountRows(db))
	requireNoError(QueryJoins(db))
	requireNoError(QueryOr(db))
	requireNoError(QueryRawSQL(db))
	requireNoError(QueryByNormalVars(db))
	requireNoError(QueryWithJsonColumn(db))
	requireNoError(UpdateByModel(db))
	requireNoError(UpdateByMap(db))
	requireNoError(DeleteById(db))
	requireNoError(Transaction(db))
	requireNoError(TransactionWrapper(db))
	requireNoError(ExecRawSQL(db))
	requireNoError(UpsertPoint(db))
	requireNoError(UpdatePointByExpress(db))
	requireNoError(DistributionLock(db))
}

func requireNoError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func requireError(err error) {
	if err == nil {
		log.Panic(err)
	}
}

/*
[单条插入]
*/
func InsertSingle(db *sqlca.Engine) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	price := float64(12.33)
	var do = models.InventoryData{
		Id:         uint64(db.NewID()),
		CreateId:   1,
		CreateName: "admin",
		CreateTime: now,
		UpdateId:   1,
		UpdateName: "admin",
		UpdateTime: now,
		IsFrozen:   models.FrozenState_Ture,
		Name:       "齿轮",
		SerialNo:   "SNO_001",
		Quantity:   1000,
		Price:      &price,
		ProductExtra: &models.ProductExtraData{
			SpecsValue: "齿数：30",
			AvgPrice:   sqlca.NewDecimal(30.8),
		},
	}

	var err error
	/*
		INSERT INTO inventory_data (`id`,`create_id`,`create_name`,`create_time`,`update_id`,`update_name`,`update_time`,`is_frozen`,`name`,`serial_no`,`quantity`,`price`,`product_extra`)
		VALUES ('1859078192380252161','1','admin','2024-11-20 11:35:55','1','admin','2024-11-20 11:35:55','0','轮胎','SNO_002','2000','210','{}')
	*/
	var rowsAffected int64
	_, rowsAffected, err = db.Model(&do).Insert()
	if err != nil {
		return log.Errorf("数据插入错误: %s", err)
	}
	log.Infof("插入数据数量：%v", rowsAffected)
	return nil
}

/*
[批量插入]
*/
func InsertBatch(db *sqlca.Engine) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	var dos = []models.InventoryData{
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
			//Price:      10.5,
			ProductExtra: nil,
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
			//Price:      210,
			ProductExtra: &models.ProductExtraData{
				SpecsValue: "17英寸",
				AvgPrice:   sqlca.NewDecimal(450.5),
			},
		},
	}

	var err error
	/*
		INSERT IGNORE INTO inventory_data
			(`id`,`create_id`,`create_name`,`create_time`,`update_id`,`update_name`,`update_time`,`is_frozen`,`name`,`serial_no`,`quantity`,`price`,`product_extra`)
		VALUES
			('1867379968636358656','1','admin','2024-12-13 09:24:13','1','admin','2024-12-13 09:24:13','0','齿轮','SNO_001','1000','10.5','{\"avg_price\":\".8\",\"specs_value\":\"齿数：32\"}'),
			('1867379968636358657','1','admin','2024-12-13 09:24:13','1','admin','2024-12-13 09:24:13','0','轮胎','SNO_002','2000','210','{\"avg_price\":\"450.5\",\"specs_value\":\"17英寸\"}')
	*/
	var rowsAffected int64
	_, rowsAffected, err = db.Model(&dos).Ignore().Insert()
	if err != nil {
		return log.Errorf("数据插入错误: %s", err)
	}
	log.Infof("批量数据插入数量：%v", rowsAffected)
	return nil
}

func UpsertSingle(db *sqlca.Engine) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	price := float64(12.33)
	var do = models.InventoryData{
		Id:         productId,
		CreateId:   1,
		CreateName: "admin",
		CreateTime: now,
		UpdateId:   1,
		UpdateName: "admin",
		UpdateTime: now,
		IsFrozen:   models.FrozenState_Ture,
		Name:       "齿轮",
		SerialNo:   "SNO_001",
		Quantity:   1000,
		Price:      &price,
		ProductExtra: &models.ProductExtraData{
			SpecsValue: "齿数：20",
			AvgPrice:   sqlca.NewDecimal(20.8),
		},
	}

	var err error
	var rowsAffected int64
	rowsAffected, err = db.Model(&do).Upsert()
	if err != nil {
		return log.Errorf("数据插入错误: %s", err)
	}
	log.Infof("插入数据数量：%v", rowsAffected)
	return nil
}

/*
[普通查询带LIMIT限制]
*/
func QueryLimit(db *sqlca.Engine) error {
	var err error
	var count int64
	var dos []*models.InventoryData
	ctx, closer := context.WithTimeout(context.Background(), 5*time.Second)
	defer closer()

	//SELECT * FROM inventory_data ORDER BY create_time DESC LIMIT 2
	count, err = db.Model(ctx, &dos).
		Select("id, name, serial_no, quantity, product_extra").
		Limit(5).
		Desc("create_time").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	log.Json("查询结果(JSON)", dos)
	return nil
}

/*
[查询无数据则报错]
*/
func QueryErrNotFound(db *sqlca.Engine) error {
	var err error
	var count int64
	var do *models.InventoryData
	ctx, closer := context.WithTimeout(context.Background(), 30*time.Second)
	defer closer()

	count, err = db.Model(ctx, &do).Id(productId).MustFind()
	if err != nil {
		if errors.Is(err, sqlca.ErrRecordNotFound) {
			log.Infof("根据ID查询数据库记录无结果：%s", err)
			return nil
		}
		return log.Errorf("数据库错误：%s", err)
	}
	log.Infof("查询结果条数: %d 数据: %+v", count, do)

	//SELECT * FROM inventory_data WHERE id=1899078192380252160
	count, err = db.Model(ctx, &do).Id(1899078192380252160).MustFind()
	if err != nil {
		if errors.Is(err, sqlca.ErrRecordNotFound) {
			log.Infof("根据ID查询数据库记录无结果：%s (正常)", err)
			return err
		}
		return log.Errorf("数据库错误：%s", err)
	}
	log.Infof("查询结果条数: %d", count)
	//log.Json("查询结果(JSON)", dos)
	return nil
}

/*
[分页查询]
*/
func QueryByPage(db *sqlca.Engine) error {
	var err error
	var count, total int64
	var dos []*models.InventoryData
	ctx, closer := context.WithTimeout(context.Background(), 30*time.Second)
	defer closer()

	//SELECT  * FROM inventory_data WHERE 1=1 ORDER BY create_time DESC LIMIT 0,20
	count, total, err = db.Model(ctx, &dos).
		Page(1, 20).
		Desc("create_time").
		QueryEx()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果条数: %d 数据库总数：%v", count, total)
	//log.Json("查询结果(JSON)", dos)
	return nil
}

/*
[复杂查询]
*/
func QueryByCondition(db *sqlca.Engine) error {
	var err error
	var count int64
	var dos []*models.InventoryData
	//SELECT * FROM inventory_data WHERE `quantity` > 0 and is_frozen IN (0,1) AND create_time >= '2024-10-01 11:35:14' ORDER BY create_time DESC
	count, err = db.Model(&dos).
		Where("is_frozen in (?)", []int{models.FrozenState_False, models.FrozenState_Ture}).
		//In("is_frozen", []int{0, 1}).
		Gt("quantity", 0).
		Gte("create_time", "2024-10-01 11:35:14").
		Desc("create_time").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}

/*
[查询带多个OR条件]
*/
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
		Limit(5).
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	//log.Json("查询结果(JSON)", dos)

	//SELECT * FROM inventory_data WHERE create_id=1 AND is_frozen = 0 AND quantity > 0 AND (name = '配件' OR serial_no = 'SNO_001') ORDER BY create_time DESC
	var andConditions = make(map[string]interface{})
	var orConditions = make(map[string]interface{})

	andConditions["create_id"] = 1    //create_id = 1
	andConditions["is_frozen"] = 0    //is_frozen = 0
	andConditions["quantity > ?"] = 0 //quantity > 0

	orConditions["name = ?"] = "配件"
	orConditions["serial_no = ?"] = "SNO_001"

	count, err = db.Model(&dos).
		And(andConditions).
		Or(orConditions).
		Desc("create_time").
		Limit(5).
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}

/*
[分组查询]
*/
func QueryByGroup(db *sqlca.Engine) error {
	var err error
	var count, total int64
	var dos []*models.InventoryData
	/*
		SELECT  create_id, SUM(quantity) AS quantity
		FROM inventory_data
		WHERE 1=1 AND quantity>'0' AND is_frozen='0' AND create_time>='2024-10-01 11:35:14'
		GROUP BY create_id
	*/
	count, total, err = db.Model(&dos).
		Select("create_id", "SUM(quantity) AS quantity").
		Gt("quantity", 0).
		Eq("is_frozen", 0).
		Gte("create_time", "2024-10-01 11:35:14").
		GroupBy("create_id", "create_name").
		QueryEx()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询返回数据条数: %d 总数：%d", count, total)
	return nil
}

// 获取查询结果行数
func QueryCountRows(db *sqlca.Engine) error {
	// SELECT COUNT(*) FROM inventory_data WHERE is_frozen = true
	count, err := db.Model(&models.InventoryData{}).Where("is_frozen", models.FrozenState_Ture).CountRows()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("select 统计总行数：%d", count)

	count, err = db.Model(&models.InventoryData{}).
		GroupBy("create_id").
		Where("create_time > ? AND is_frozen = ?", "2025-06-01 00:00:00", models.FrozenState_False).
		CountRows()
	log.Infof("group by 统计总行数：%d", count)
	return nil
}

/*
[联表查询]
*/
func QueryJoins(db *sqlca.Engine) error {
	/*
		SELECT a.id as product_id, a.name AS product_name, b.quantity, b.weight
		FROM inventory_data a
		LEFT JOIN inventory_in b
		ON a.id=b.product_id
		WHERE a.quantity > 0 AND a.is_frozen=0 AND a.create_time>='2024-10-01 11:35:14'
	*/
	var do struct{}
	count, err := db.Model(&do).
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

/*
[普通变量取值查询]
*/

func QueryByNormalVars(db *sqlca.Engine) error {
	var err error
	var name, serialNo string
	var id = uint64(productId)
	//SELECT name, serial_no FROM inventory_data WHERE id=1906626367382884352
	_, err = db.Model(&name, &serialNo).
		Table("inventory_data").
		Select("name, serial_no").
		Id(id).
		MustFind()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("数据ID: %v name=%s serial_no=%s", id, name, serialNo)
	return nil
}

/*
[查询保存JSON内容的字段到结构体]
models.InventoryData对象的ProductExtra是一个跟数据库JSON内容对应的结构体
*/
func QueryWithJsonColumn(db *sqlca.Engine) error {
	var err error
	var do models.InventoryData
	var id = uint64(productId)

	/*
		SELECT * FROM inventory_data WHERE id=1906626367382884352

		+-----------------------+-----------------------+-----------------------+------------------------------------------------+
		| id	                | name	| serial_no	    | quantity	| price	    |                 product_extra                  |
		+-----------------------+-------+---------------+-----------+-----------+------------------------------------------------+
		| 1906626367382884352	| 轮胎  	| SNO_002		| 2000.000 	| 210.00	| {"avg_price": "450.5", "specs_value": "17英寸"} |
		+------------------------------------------------------------------------------------------------------------------------+
	*/
	_, err = db.Model(&do).
		Table("inventory_data").
		Select("id", "name", "serial_no", "quantity", "price", "product_extra").
		Id(id).
		MustFind()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Json(do)
	/*
		2024-12-18 15:15:03.560732 PID:64764 [INFO] {goroutine 1} <main.go:373 QueryWithJsonColumn()> ID: 1906626367382884352 数据：{Id:1906626367382884352 Name:轮胎 SerialNo:SNO_002 Quantity:2000 Price:210 ProductExtra:{AvgPrice:450.5 SpecsValue:17英寸}}
	*/
	return nil
}

/*
[常规SQL查询]
*/
func QueryRawSQL(db *sqlca.Engine) error {
	var rows []*models.InventoryData

	//SELECT * FROM inventory_data  WHERE is_frozen IN (1) AND quantity > '10'
	_, err := db.Model(&rows).QueryRaw("SELECT * FROM inventory_data WHERE is_frozen in (?) AND quantity > ?", []models.FrozenState{0, 1}, 10)
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	return nil
}

func ExecRawSQL(db *sqlca.Engine) error {
	var sb = sqlca.NewStringBuilder()

	//UPDATE inventory_data SET quantity = '10' WHERE id=1867379968636358657
	sb.Append("UPDATE inventory_data")
	sb.Append("SET quantity = ?", 10)
	sb.Append("WHERE id = ?", productId)

	strQuery := sb.String()
	affectedRows, lastInsertId, err := db.Model(nil).ExecRaw(strQuery)
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("受影响的行数：%d 最后插入的ID：%d", affectedRows, lastInsertId)
	return nil
}

/*
[数据更新]

SELECT * FROM inventory_data  WHERE `id`='1906626367382884352'
UPDATE inventory_data SET `quantity`='2300' WHERE `id`='1906626367382884352'
*/
func UpdateByModel(db *sqlca.Engine) error {
	var err error
	var do *models.InventoryData
	var id = uint64(productId)
	_, err = db.Model(&do).Id(id).MustFind() //Find方法如果是单条记录没找到则提示ErrNotFound错误（Query方法不会报错）
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}

	do.Quantity = 2300 //更改库存
	_, err = db.Model(&do).Select("quantity").Update()
	if err != nil {
		return log.Errorf("更新错误：%s", err)
	}
	return nil
}

/*
[通过map进行数据更新]
*/
func UpdateByMap(db *sqlca.Engine) error {
	var err error
	var id = uint64(productId)
	var updates = map[string]interface{}{
		"quantity":      2100, //更改库存
		"price":         300,  //更改价格
		"product_extra": nil,  //设置产品扩展数据为NULL
	}
	//UPDATE inventory_data SET `quantity`='2100',`price`=300, is_frozen = NULL WHERE `id`='1906626367382884352'
	_, err = db.Model(updates).Table("inventory_data").Id(id).Update()
	if err != nil {
		return log.Errorf("更新错误：%s", err)
	}
	return nil
}

/*
[删除数据]
*/
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

/*
[事务处理]
*/
func Transaction(db *sqlca.Engine) error {

	/*
		-- TRANSACTION BEGIN

			INSERT INTO inventory_in (`user_id`,`quantity`,`remark`,`create_id`,`user_name`,`weight`,`create_time`,`update_name`,`is_deleted`,`product_id`,`id`,`create_name`,`update_id`,`update_time`,`order_no`) VALUES ('3','20','产品入库','1','lazy','200.3','2024-11-27 11:35:14','admin','0','1906626367382884352','1861614736295071744','admin','1','2024-11-27 1114','202407090000001')
			SELECT * FROM inventory_data  WHERE `id`='1906626367382884352'
			UPDATE inventory_data SET `quantity`='2320' WHERE `id`='1906626367382884352'

		-- TRANSACTION END
	*/

	now := time.Now().Format("2006-01-02 15:04:05")
	tx, err := db.TxBegin()
	if err != nil {
		return log.Errorf("开启事务失败：%s", err)
	}
	defer tx.TxRollback()

	productId := uint64(productId)
	strOrderNo := time.Now().Format("20060102150405.000000000")
	//***************** 执行事务操作 *****************
	quantity := float64(20)
	weight := float64(200.3)
	_, _, err = tx.Model(&models.InventoryIn{
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
	_, err = tx.Model(&inventoryData).Id(productId).MustFind() //Find方法如果是单条记录没找到则提示ErrNotFound错误（Query方法不会报错）
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

/*
[事务处理封装]
*/
func TransactionWrapper(db *sqlca.Engine) error {
	/*
	   -- TRANSACTION BEGIN

	   	INSERT INTO inventory_in (`user_id`,`quantity`,`remark`,`create_id`,`user_name`,`weight`,`create_time`,`update_name`,`is_deleted`,`product_id`,`id`,`create_name`,`update_id`,`update_time`,`order_no`) VALUES ('3','20','产品入库','1','lazy','200.3','2024-11-27 11:35:14','admin','0','1906626367382884352','1861614736295071744','admin','1','2024-11-27 1114','202407090000002')
	   	SELECT * FROM inventory_data  WHERE `id`='1906626367382884352'
	   	UPDATE inventory_data SET `quantity`='2320' WHERE `id`='1906626367382884352'

	   -- TRANSACTION END
	*/
	strOrderNo := time.Now().Format("20060102150405.000000000")
	err := db.TxFunc(func(tx *sqlca.Engine) error {
		var err error
		now := time.Now().Format("2006-01-02 15:04:05")

		//***************** 执行事务操作 *****************
		quantity := float64(20)
		weight := float64(200.3)
		_, _, err = tx.Model(&models.InventoryIn{
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
		_, err = tx.Model(&inventoryData).Id(productId).MustFind() //Find方法如果是单条记录没找到则提示ErrNotFound错误（Query方法不会报错）
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

// 地理位置坐标插入/更新
func UpsertPoint(db *sqlca.Engine) error {
	now := time.Now().Format(time.DateTime)
	price := 243.3
	id := db.NewID()
	do := &models.InventoryData{
		Id:         uint64(id),
		CreateId:   1,
		CreateName: "admin",
		CreateTime: now,
		UpdateId:   1,
		UpdateName: "admin",
		UpdateTime: now,
		IsFrozen:   models.FrozenState_Ture,
		Name:       "齿轮",
		SerialNo:   "SNO_001",
		Quantity:   1000,
		Price:      &price,
		Location: sqlca.Point{
			X: 112.34232,
			Y: -20.32432,
		},
	}
	_, _, err := db.Model(&do).Insert()
	if err != nil {
		return log.Errorf(err.Error())
	}
	do.Location.X = 111.23
	do.Location.Y = -21.53
	_, err = db.Model(&do).Select("location").Update()
	if err != nil {
		return log.Errorf(err.Error())
	}
	var updates = map[string]any{
		"location": sqlca.Point{X: 110.234, Y: -10.23},
	}
	_, err = db.Model(updates).Table("inventory_data").Id(id).Update()
	if err != nil {
		return log.Errorf(err.Error())
	}

	var do2 *models.InventoryData
	_, err = db.Model(&do2).Id(id).Query()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("do2 %+v", do2)
	return nil
}

func UpdatePointByExpress(db *sqlca.Engine) error {
	var upmap = map[string]any{
		"location": sqlca.NewExpr("POINT(?, ?)", 113.234, 22.39236),
	}
	rows, err := db.Model(upmap).Table("inventory_data").Id(productId).Update()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("rows affected: %d", rows)
	return nil
}

func DistributionLock(db *sqlca.Engine) error {
	key := fmt.Sprintf("test:inventory_data:product_id:%v", productId)
	unlock, err := db.Lock(key, 10*time.Second)
	if err != nil {
		return log.Errorf(err.Error())
	}
	defer unlock()

	var upmap = map[string]any{
		"location": sqlca.NewExpr("POINT(?, ?)", 116.2, 22.1),
	}
	rows, err := db.Model(upmap).Table("inventory_data").Id(productId).Update()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("rows affected: %d", rows)
	return nil
}
