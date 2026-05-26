package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3"
	"github.com/civet148/sqlca/v3/example/models"
)

var db *sqlca.Engine

func init() {
	var err error
	var opts = []sqlca.Option{
		sqlca.WithDebug(),
		sqlca.WithMaxConn(100),
		sqlca.WithIdleConn(5),
		//sqlca.WithSSH(&sqlca.SSH{ //SSH tunnel config
		//	User:     "root",
		//	Password: "123456",
		//	Host:     "192.168.0.19:22",
		//}),
		sqlca.WithRedisConfig(&sqlca.RedisConfig{ //redis distribution lock config
			Address: "127.0.0.1:6379",
		}),
		sqlca.WithSnowFlake(&sqlca.SnowFlake{ // snowflake algorithm config
			NodeId: 1,
		}),
		sqlca.WithAutoMigrate(),
	}
	db, err = sqlca.NewEngine("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4", opts...)
	if err != nil {
		log.Errorf("connect database error: %s", err)
		return
	}
	// 1. 自动迁移模型
	requireNoError(MigrateAllModels(db))

	// 2. 清理所有数据，确保数据库为空
	requireNoError(CleanAllData(db))

	// 3. 插入测试数据
	requireNoError(InsertTestData(db))

	// 4. 验证数据插入
	requireNoError(VerifyTestData(db))
}

func main() {
	// 测试普通变量查询
	requireNoError(TestQueryByNormalVars(db))

	// 测试关联查询
	requireNoError(TestPreload(db))

	// 测试更新功能
	requireNoError(TestUpdate(db))

	// 测试删除功能
	requireNoError(TestDelete(db))

	// 测试事务处理
	requireNoError(TestTransaction(db))

	// 测试批量插入
	requireNoError(TestInsertBatch(db))

	// 测试查询限制
	requireNoError(TestQueryLimit(db))

	// 测试查询无结果
	requireError(TestQueryErrNotFound(db))

	// 测试分页查询
	requireNoError(TestQueryByPage(db))

	// 测试条件查询
	requireNoError(TestQueryByCondition(db))

	// 测试分组查询
	requireNoError(TestQueryByGroup(db))

	// 测试统计行数
	requireNoError(TestQueryCountRows(db))

	// 测试联表查询
	requireNoError(TestQueryJoins(db))

	// 测试 OR 条件查询
	requireNoError(TestQueryOr(db))

	// 测试原始 SQL 查询
	requireNoError(TestQueryRawSQL(db))

	// 测试 JSON 字段查询
	requireNoError(TestQueryWithJsonColumn(db))

	// 测试通过 Map 更新
	requireNoError(TestUpdateByMap(db))

	// 测试事务封装
	requireNoError(TestTransactionWrapper(db))

	// 测试执行原始 SQL
	requireNoError(TestExecRawSQL(db))

	// 测试地理位置坐标操作
	requireNoError(TestUpsertPoint(db))
	requireNoError(TestUpdatePointByExpress(db))

	// 测试分布式锁
	// requireNoError(TestDistributionLock(db))

	log.Infof("所有测试验证通过！")
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

func MigrateAllModels(db *sqlca.Engine) (err error) {
	return db.AutoMigrate(context.Background(), nil,
		&models.User{}, &models.UserProfile{}, &models.Role{}, &models.UserRole{},
		&models.InventoryData{}, &models.InventoryIn{}, &models.InventoryOut{})
}

func CleanAllData(db *sqlca.Engine) (err error) {
	// 按依赖关系顺序删除数据
	tables := []string{"user_roles", "user_profiles", "roles", "users", "inventory_in", "inventory_out", "inventory_data"}
	for _, table := range tables {
		_, _, err = db.Model(nil).ExecRaw(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return log.Errorf("清理表 %s 数据失败: %s", table, err)
		}
	}
	log.Infof("所有表数据清理完成")
	return nil
}

func InsertTestData(db *sqlca.Engine) (err error) {
	// 1. 插入角色数据
	roles := []*models.Role{
		{
			Id: 1,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "admin",
		},
		{
			Id: 2,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "user",
		},
	}

	for _, role := range roles {
		_, err = db.Model(role).Upsert()
		if err != nil {
			return log.Errorf("插入角色数据失败: %s", err)
		}
	}

	// 2. 插入用户数据
	users := []*models.User{
		{
			Id: 1,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserName: "lory",
			Email:    "lory@hotmail.com",
		},
		{
			Id: 2,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserName: "civet148",
			Email:    "civet148@126.com",
		},
	}

	for _, user := range users {
		_, err = db.Model(user).Upsert()
		if err != nil {
			return log.Errorf("插入用户数据失败: %s", err)
		}
	}

	// 3. 插入用户资料数据
	profiles := []*models.UserProfile{
		{
			Id: 1,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:  1,
			Avatar:  "https://www.hello.com/lory.jpg",
			Address: "中国北京市朝阳区A座",
		},
		{
			Id: 2,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:  2,
			Avatar:  "https://www.hello.com/civet148.jpg",
			Address: "中国上海市浦东新区B座",
		},
	}

	for _, profile := range profiles {
		_, err = db.Model(profile).Upsert()
		if err != nil {
			return log.Errorf("插入用户资料数据失败: %s", err)
		}
	}

	// 4. 插入用户角色关联数据
	userRoles := []*models.UserRole{
		{
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId: 1,
			RoleId: 1,
		},
		{
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId: 1,
			RoleId: 2,
		},
		{
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId: 2,
			RoleId: 2,
		},
	}

	for _, userRole := range userRoles {
		_, err = db.Model(userRole).Upsert()
		if err != nil {
			return log.Errorf("插入用户角色关联数据失败: %s", err)
		}
	}

	// 5. 插入库存数据
	price := float64(12.33)
	inventoryData := &models.InventoryData{
		Id: uint64(db.NewID()),
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		IsFrozen: models.FrozenState_False,
		Name:     "齿轮",
		SerialNo: "SNO_001",
		Quantity: 1000,
		Price:    &price,
		ProductExtra: &models.ProductExtraData{
			SpecsValue: "齿数：30",
			AvgPrice:   sqlca.NewDecimal(30.8),
		},
	}

	_, _, err = db.Model(inventoryData).Insert()
	if err != nil {
		return log.Errorf("插入库存数据失败: %s", err)
	}

	log.Infof("测试数据插入完成")
	return nil
}

func VerifyTestData(db *sqlca.Engine) (err error) {
	// 1. 验证角色数据
	var roles []*models.Role
	count, err := db.Model(&roles).Query()
	if err != nil {
		return log.Errorf("查询角色数据失败: %s", err)
	}
	if count != 2 {
		return log.Errorf("角色数据数量不正确，期望 2，实际 %d", count)
	}
	log.Infof("角色数据验证通过，共 %d 条", count)

	// 2. 验证用户数据
	var users []*models.User
	count, err = db.Model(&users).Query()
	if err != nil {
		return log.Errorf("查询用户数据失败: %s", err)
	}
	if count != 2 {
		return log.Errorf("用户数据数量不正确，期望 2，实际 %d", count)
	}
	log.Infof("用户数据验证通过，共 %d 条", count)

	// 3. 验证用户资料数据
	var profiles []*models.UserProfile
	count, err = db.Model(&profiles).Query()
	if err != nil {
		return log.Errorf("查询用户资料数据失败: %s", err)
	}
	if count != 2 {
		return log.Errorf("用户资料数据数量不正确，期望 2，实际 %d", count)
	}
	log.Infof("用户资料数据验证通过，共 %d 条", count)

	// 4. 验证库存数据
	var inventoryData []*models.InventoryData
	count, err = db.Model(&inventoryData).Query()
	if err != nil {
		return log.Errorf("查询库存数据失败: %s", err)
	}
	if count != 1 {
		return log.Errorf("库存数据数量不正确，期望 1，实际 %d", count)
	}
	log.Infof("库存数据验证通过，共 %d 条", count)

	return nil
}

func TestPreload(db *sqlca.Engine) (err error) {
	// 测试关联查询
	var users []*models.User
	rows, err := db.Model(&users).Preload("Roles", "id > ?", 0).Preload("Profile").Query()
	if err != nil {
		return log.Errorf("关联查询失败: %s", err)
	}

	if len(users) != 2 {
		return log.Errorf("关联查询用户数量不正确，期望 2，实际 %d", len(users))
	}

	// 验证第一个用户的关联数据
	user1 := users[0]
	if len(user1.Roles) != 2 {
		return log.Errorf("用户 %s 的角色数量不正确，期望 2，实际 %d", user1.UserName, len(user1.Roles))
	}
	if user1.Profile.UserId != user1.Id {
		return log.Errorf("用户 %s 的资料关联不正确", user1.UserName)
	}
	log.Json("user1", user1)

	// 验证第二个用户的关联数据
	user2 := users[1]
	if len(user2.Roles) != 1 {
		return log.Errorf("用户 %s 的角色数量不正确，期望 1，实际 %d", user2.UserName, len(user2.Roles))
	}
	if user2.Profile.UserId != user2.Id {
		return log.Errorf("用户 %s 的资料关联不正确", user2.UserName)
	}
	log.Json("user2", user2)
	log.Infof("关联查询测试通过，共 %d 条记录", rows)
	return nil
}

func TestUpdate(db *sqlca.Engine) (err error) {
	// 1. 更新用户信息
	var user *models.User
	_, err = db.Model(&user).Id(1).MustFind()
	if err != nil {
		return log.Errorf("查询用户失败: %s", err)
	}

	oldEmail := user.Email
	user.Email = "lory_updated@hotmail.com"
	_, err = db.Model(&user).Select("email").Update()
	if err != nil {
		return log.Errorf("更新用户失败: %s", err)
	}

	// 验证更新结果
	var updatedUser *models.User
	_, err = db.Model(&updatedUser).Id(1).MustFind()
	if err != nil {
		return log.Errorf("查询更新后的用户失败: %s", err)
	}

	if updatedUser.Email == oldEmail {
		return log.Errorf("用户邮箱未更新")
	}

	// 2. 更新库存信息
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	oldQuantity := inventory.Quantity
	inventory.Quantity = 1500
	_, err = db.Model(&inventory).Select("quantity").Update()
	if err != nil {
		return log.Errorf("更新库存失败: %s", err)
	}

	// 验证更新结果
	var updatedInventory *models.InventoryData
	_, err = db.Model(&updatedInventory).Id(inventory.Id).MustFind()
	if err != nil {
		return log.Errorf("查询更新后的库存失败: %s", err)
	}

	if updatedInventory.Quantity == oldQuantity {
		return log.Errorf("库存数量未更新")
	}

	log.Infof("更新功能测试通过")
	return nil
}

func TestDelete(db *sqlca.Engine) (err error) {
	// 1. 删除用户角色关联
	_, err = db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", 1, 1).Delete()
	if err != nil {
		return log.Errorf("删除用户角色关联失败: %s", err)
	}

	// 验证删除结果
	var userRoles []*models.UserRole
	count, err := db.Model(&userRoles).Where("user_id = ?", 1).Query()
	if err != nil {
		return log.Errorf("查询用户角色关联失败: %s", err)
	}

	if count != 1 {
		return log.Errorf("用户角色关联删除失败，期望 1，实际 %d", count)
	}

	// 2. 删除库存数据
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	inventoryId := inventory.Id
	_, err = db.Model(&models.InventoryData{}).Id(inventoryId).Delete()
	if err != nil {
		return log.Errorf("删除库存失败: %s", err)
	}

	// 验证删除结果
	var deletedInventory *models.InventoryData
	_, err = db.Model(&deletedInventory).Id(inventoryId).MustFind()
	if err == nil {
		return log.Errorf("库存数据未删除")
	}
	// 确认错误是因为记录未找到
	if !strings.Contains(err.Error(), "record not found") {
		return log.Errorf("查询删除后的库存失败: %s", err)
	}

	log.Infof("删除功能测试通过")
	return nil
}

func TestTransaction(db *sqlca.Engine) (err error) {
	// 测试事务处理
	tx, err := db.TxBegin()
	if err != nil {
		return log.Errorf("开启事务失败: %s", err)
	}
	defer tx.TxRollback()

	// 1. 在事务中插入新用户
	newUser := &models.User{
		Id: 3,
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName: "testuser",
		Email:    "test@example.com",
	}

	_, err = tx.Model(newUser).Upsert()
	if err != nil {
		return log.Errorf("事务中插入用户失败: %s", err)
	}

	// 2. 在事务中插入用户资料
	newProfile := &models.UserProfile{
		Id: 3,
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserId:  3,
		Avatar:  "https://www.hello.com/test.jpg",
		Address: "中国广州市天河区C座",
	}

	_, err = tx.Model(newProfile).Upsert()
	if err != nil {
		return log.Errorf("事务中插入用户资料失败: %s", err)
	}

	// 3. 提交事务
	err = tx.TxCommit()
	if err != nil {
		return log.Errorf("提交事务失败: %s", err)
	}

	// 验证事务结果
	var user *models.User
	_, err = db.Model(&user).Id(3).Preload("Profile").Query()
	if err != nil {
		return log.Errorf("查询事务中创建的用户失败: %s", err)
	}

	if user == nil {
		return log.Errorf("事务中创建的用户未找到")
	}

	if user.Profile.UserId != user.Id {
		return log.Errorf("事务中创建的用户资料关联不正确")
	}

	log.Infof("事务处理测试通过")
	return nil
}

func TestInsertBatch(db *sqlca.Engine) (err error) {
	// 批量插入测试
	var dos = []models.InventoryData{
		{
			Id: uint64(db.NewID()),
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			IsFrozen:     models.FrozenState_False,
			Name:         "齿轮",
			SerialNo:     "SNO_002",
			Quantity:     1000,
			ProductExtra: nil,
		},
		{
			Id: uint64(db.NewID()),
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			IsFrozen: models.FrozenState_False,
			Name:     "轮胎",
			SerialNo: "SNO_003",
			Quantity: 2000,
			ProductExtra: &models.ProductExtraData{
				SpecsValue: "17英寸",
				AvgPrice:   sqlca.NewDecimal(450.5),
			},
		},
	}

	_, rowsAffected, err := db.Model(&dos).Ignore().Insert()
	if err != nil {
		return log.Errorf("批量数据插入错误: %s", err)
	}
	log.Infof("批量数据插入数量：%v", rowsAffected)
	return nil
}

func TestQueryLimit(db *sqlca.Engine) (err error) {
	// 测试查询限制
	var dos []*models.InventoryData
	count, err := db.Model(&dos).
		Select("id, name, serial_no, quantity, product_extra").
		Limit(5).
		Desc("created_at").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("1.查询结果数据条数: %d", count)
	return nil
}

func TestQueryErrNotFound(db *sqlca.Engine) (err error) {
	// 测试查询无结果
	var do *models.InventoryData
	_, err = db.Model(&do).Id(uint64(9999999999999999999)).MustFind()
	if err == nil {
		return log.Errorf("应该返回错误，但没有返回")
	}
	log.Infof("查询无结果测试通过: %s", err)
	return err
}

func TestQueryByPage(db *sqlca.Engine) (err error) {
	// 测试分页查询
	var dos []*models.InventoryData
	count, total, err := db.Model(&dos).
		Page(1, 20).
		Desc("created_at").
		QueryEx()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果条数: %d 数据库总数：%v", count, total)
	return nil
}

func TestQueryByCondition(db *sqlca.Engine) (err error) {
	// 测试条件查询
	var dos []*models.InventoryData
	count, err := db.Model(&dos).
		Where("is_frozen in (?)", []int{models.FrozenState_False, models.FrozenState_Ture}).
		Gt("quantity", 0).
		Gte("created_at", "2026-01-01 00:00:00").
		Desc("created_at").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}

func TestQueryByGroup(db *sqlca.Engine) (err error) {
	// 测试分组查询
	var dos []*models.InventoryData
	count, total, err := db.Model(&dos).
		Select("is_frozen", "SUM(quantity) AS quantity").
		Gt("quantity", 0).
		GroupBy("is_frozen").
		QueryEx()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询返回数据条数: %d 总数：%d", count, total)
	return nil
}

func TestQueryCountRows(db *sqlca.Engine) (err error) {
	// 测试统计行数
	count, err := db.Model(&models.InventoryData{}).Where("is_frozen", models.FrozenState_False).CountRows()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("select 统计总行数：%d", count)

	count, err = db.Model(&models.InventoryData{}).
		GroupBy("is_frozen").
		Where("created_at > ? AND is_frozen = ?", "2026-01-01 00:00:00", models.FrozenState_False).
		CountRows()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("group by 统计总行数：%d", count)
	return nil
}

func TestQueryJoins(db *sqlca.Engine) (err error) {
	// 测试联表查询
	var do struct{}
	count, err := db.Model(&do).
		Select("a.id as product_id", "a.name AS product_name", "b.quantity", "b.weight").
		Table("inventory_data a").
		LeftJoin("inventory_in b").
		On("a.id=b.product_id").
		Gt("a.quantity", 0).
		Eq("a.is_frozen", 0).
		Gte("a.created_at", "2026-01-01 00:00:00").
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)
	return nil
}

func TestQueryOr(db *sqlca.Engine) (err error) {
	// 测试 OR 条件查询
	var dos []*models.InventoryData

	count, err := db.Model(&dos).
		And("create_id = ?", 0).
		Or("name = ?", "齿轮").
		Or("serial_no = ?", "SNO_001").
		Desc("created_at").
		Limit(5).
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("查询结果数据条数: %d", count)

	// 测试 AND 和 OR 组合条件
	var andConditions = make(map[string]interface{})
	var orConditions = make(map[string]interface{})

	andConditions["create_id"] = 0    //create_id = 0
	andConditions["is_frozen"] = 0    //is_frozen = 0
	andConditions["quantity > ?"] = 0 //quantity > 0

	orConditions["name = ?"] = "齿轮"
	orConditions["serial_no = ?"] = "SNO_001"

	count, err = db.Model(&dos).
		And(andConditions).
		Or(orConditions).
		Desc("created_at").
		Limit(5).
		Query()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("组合条件查询结果数据条数: %d", count)
	return nil
}

func TestQueryRawSQL(db *sqlca.Engine) (err error) {
	// 测试原始 SQL 查询
	var rows []*models.InventoryData

	_, err = db.Model(&rows).QueryRaw("SELECT * FROM inventory_data WHERE is_frozen in (?) AND quantity > ?", []models.FrozenState{0, 1}, 10)
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("原始 SQL 查询成功，返回 %d 条记录", len(rows))
	return nil
}

func TestQueryByNormalVars(db *sqlca.Engine) (err error) {
	// 测试普通变量查询
	var name, serialNo string
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	_, err = db.Model(&name, &serialNo).
		Debug().
		Table("inventory_data").
		Select("name, serial_no").
		Id(inventory.Id).
		MustFind()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("数据ID: %v name=%s serial_no=%s", inventory.Id, name, serialNo)
	return nil
}

func TestQueryWithJsonColumn(db *sqlca.Engine) (err error) {
	// 测试 JSON 字段查询
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	var do models.InventoryData
	_, err = db.Model(&do).
		Table("inventory_data").
		Select("id", "name", "serial_no", "quantity", "price", "product_extra").
		Id(inventory.Id).
		MustFind()
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Json(do)
	return nil
}

func TestUpdateByMap(db *sqlca.Engine) (err error) {
	// 测试通过 Map 更新
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	var updates = map[string]interface{}{
		"quantity":      2100, //更改库存
		"price":         300,  //更改价格
		"product_extra": nil,  //设置产品扩展数据为NULL
	}

	_, err = db.Model(updates).Table("inventory_data").Id(inventory.Id).Update()
	if err != nil {
		return log.Errorf("更新错误：%s", err)
	}
	log.Infof("通过 Map 更新成功")
	return nil
}

func TestTransactionWrapper(db *sqlca.Engine) (err error) {
	// 测试事务封装
	strOrderNo := time.Now().Format("20060102150405.000000000")
	err = db.TxFunc(func(tx *sqlca.Engine) error {
		var err error
		// 执行事务操作
		quantity := float64(20)
		weight := float64(200.3)
		_, _, err = tx.Model(&models.InventoryIn{
			Id: uint64(db.NewID()),
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ProductId: 1,
			OrderNo:   strOrderNo,
			UserId:    3,
			UserName:  "testuser",
			Quantity:  quantity,
			Weight:    sqlca.NewDecimal(weight),
			Remark:    "产品入库",
		}).Insert()
		if err != nil {
			return log.Errorf("数据插入错误: %s", err)
		}
		return nil
	})

	if err != nil {
		return log.Errorf("事务失败：%s", err)
	}
	log.Infof("事务封装测试通过")
	return nil
}

func TestExecRawSQL(db *sqlca.Engine) (err error) {
	// 测试执行原始 SQL
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	var sb = sqlca.NewStringBuilder()
	sb.Append("UPDATE inventory_data")
	sb.Append("SET quantity = ?", 10)
	sb.Append("WHERE id = ?", inventory.Id)

	strQuery := sb.String()
	affectedRows, lastInsertId, err := db.Model(nil).ExecRaw(strQuery)
	if err != nil {
		return log.Errorf("数据查询错误：%s", err)
	}
	log.Infof("受影响的行数：%d 最后插入的ID：%d", affectedRows, lastInsertId)
	return nil
}

func TestUpsertPoint(db *sqlca.Engine) (err error) {
	// 测试地理位置坐标插入/更新
	price := 243.3
	do := &models.InventoryData{
		Id: uint64(db.NewID()),
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		IsFrozen: models.FrozenState_Ture,
		Name:     "测试坐标",
		SerialNo: "SNO_004",
		Quantity: 1000,
		Price:    &price,
		Location: sqlca.Point{
			X: 112.34232,
			Y: -20.32432,
		},
	}
	_, _, err = db.Model(&do).Insert()
	if err != nil {
		return log.Errorf(err.Error())
	}
	do.Location.X = 111.23
	do.Location.Y = -21.53
	_, err = db.Model(&do).Select("location").Update()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("地理位置坐标操作测试通过")
	return nil
}

func TestUpdatePointByExpress(db *sqlca.Engine) (err error) {
	// 测试通过表达式更新地理位置坐标
	var inventory *models.InventoryData
	_, err = db.Model(&inventory).Query()
	if err != nil {
		return log.Errorf("查询库存失败: %s", err)
	}

	var upmap = map[string]any{
		"location": sqlca.NewExpr("POINT(?, ?)", 113.234, 22.39236),
	}
	rows, err := db.Model(upmap).Table("inventory_data").Id(inventory.Id).Update()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("通过表达式更新地理位置坐标，受影响的行数：%d", rows)
	return nil
}

func TestDistributionLock(db *sqlca.Engine) (err error) {
	// 测试分布式锁
	key := fmt.Sprintf("test:inventory_data:product_id:%v", 1)
	unlock, err := db.Lock(key, 10*time.Second)
	if err != nil {
		return log.Errorf(err.Error())
	}
	defer unlock()

	var upmap = map[string]any{
		"location": sqlca.NewExpr("POINT(?, ?)", 116.2, 22.1),
	}
	rows, err := db.Model(upmap).Table("inventory_data").Id(1).Update()
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("分布式锁测试通过，受影响的行数：%d", rows)
	return nil
}
