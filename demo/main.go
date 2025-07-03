package main

import (
	"errors"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3"
	"github.com/civet148/sqlca/v3/demo/models"
	"time"
)

const (
	exampleId = 1939612151790440448
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
		//SSH: &sqlca.SSH{ //SSH隧道连接配置
		//	User:     "root",
		//	Password: "123456",
		//	Host:     "192.168.2.19:22",
		//},
	}
	db, err = sqlca.NewEngine("mysql://root:12345678@127.0.0.1:3306/test?charset=utf8mb4", options)
	if err != nil {
		log.Errorf("connect database error: %s", err)
		return
	}

	requireNoError(InsertSingle(db))
	requireNoError(InsertBatch(db))
	requireNoError(QueryLimit(db))
	requireError(QueryErrNotFound(db))
	requireNoError(QueryByPage(db))
	requireNoError(QueryByCondition(db))
	requireNoError(QueryByGroup(db))
	requireNoError(QueryJoins(db))
	requireNoError(QueryByNormalVars(db))
	requireNoError(UpdateByModel(db))
	requireNoError(UpdateByMap(db))
	requireNoError(DeleteById(db))
	requireNoError(Transaction(db))
	requireNoError(TransactionWrapper(db))
	requireNoError(QueryOr(db))
	requireNoError(QueryRawSQL(db))
	requireNoError(ExecRawSQL(db))
	requireNoError(QueryWithJsonColumn(db))
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
		IsFrozen:   0,
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
	_, err = db.Model(&do).Insert()
	if err != nil {
		return log.Errorf("数据插入错误: %s", err)
	}
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
			ProductExtra: &models.ProductExtraData{
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
			//Price:      210,
			ProductExtra: &models.ProductExtraData{
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

/*
[普通查询带LIMIT限制]
*/
func QueryLimit(db *sqlca.Engine) error {
	var err error
	var count int64
	var dos []*models.InventoryData

	//SELECT * FROM inventory_data ORDER BY create_time DESC LIMIT 2
	count, err = db.Model(&dos).
		Select("id, name, serial_no, quantity").
		Limit(2).
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

	count, err = db.Model(&do).Id(exampleId).MustFind()
	if err != nil {
		if errors.Is(err, sqlca.ErrRecordNotFound) {
			log.Infof("根据ID查询数据库记录无结果：%s", err)
			return nil
		}
		return log.Errorf("数据库错误：%s", err)
	}
	log.Infof("查询结果条数: %d 数据: %+v", count, do)

	//SELECT * FROM inventory_data WHERE id=1899078192380252160
	count, err = db.Model(&do).Id(1899078192380252160).MustFind()
	if err != nil {
		if errors.Is(err, sqlca.ErrRecordNotFound) {
			return log.Errorf("根据ID查询数据库记录无结果：%s", err)
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

	//SELECT  * FROM inventory_data WHERE 1=1 ORDER BY create_time DESC LIMIT 0,20
	count, total, err = db.Model(&dos).
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
	var id = uint64(exampleId)
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
	var id = uint64(exampleId)

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

/*
[数据更新]

SELECT * FROM inventory_data  WHERE `id`='1906626367382884352'
UPDATE inventory_data SET `quantity`='2300' WHERE `id`='1906626367382884352'
*/
func UpdateByModel(db *sqlca.Engine) error {
	var err error
	var do *models.InventoryData
	var id = uint64(exampleId)
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
	var id = uint64(exampleId)
	var updates = map[string]interface{}{
		"quantity": 2100, //更改库存
		"Price":    300,  //更改价格
	}
	//UPDATE inventory_data SET `quantity`='2100',`price`=300 WHERE `id`='1906626367382884352'
	_, err = db.Model(&updates).Table("inventory_data").Id(id).Update()
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

	productId := uint64(1906626367382884352)
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
		productId := uint64(exampleId)
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
