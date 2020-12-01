# author 
lory.li
# email
civet148@126.com
# QQ 
93864947
# sqlca
a enhancement database and cache tool based on sqlx and redigogo which based on redigo and go-redis-cluster

# 功能特性
* 支持数据库多主多从，数据库支持mysql/postgres/mssql
* 支持通过结构体切片、基本类型变量获取查询数据。使用基本变量时仅取第一行数据字段（按顺序赋值）
* 支持插入时如有主键或唯一索引冲突改为更新（upsert）
* 支持结构体成员自定义tag（内建支持db、protobuf、json）
* 支持基本事务和封装事务
* 支持通过主键和索引缓存到redis（仅orm操作有效）
* 命令行工具db2go支持生成数据表结构到.go文件和.proto文件
* 支持通过sqlca标签指定某些字段为readonly，指定readonly的字段在orm插入和更新时将被忽略
* 支持INNER/LEFT/RIGHT JOIN
* 支持传入结构体指针地址引用自动分配内存并赋值
* 支持通过数据对象嵌套子对象将json/text字段中的JSON内容反序列化到子对象中(或子对象切片)
* 支持GEO HASH生成
* 支持'附近'经纬度查询
* 支持数据对象模型切片进行批量插入
* 支持自定义sql.Scanner实现(赋值)
* 支持Case...When语法
* 支持SSH tunnel

# 测试案例

[数据库表结构](/test/test.sql)

[测试代码](test/main.go)
