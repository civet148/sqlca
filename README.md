# author 
lory.li
# email
civet148@126.com
# QQ 
93864947
# sqlca
a enhancement database and cache tool based on sqlx and redigogo which based on redigo and go-redis-cluster

# 中文
[中文文档](README_CN.md)

# Overview

- Almost full-featured ORM
- Multiple databases (MySQL/Postgres/MS-SQLSERVER), read/write splitting
- Multiple model type `struct, slice, built-in types, map` 
- Batch insert
- Query/Insert/Update
- Upsert by customization [only MySQL/Postgres]
- Transactions wrapper (auto rollback or commit)
- Slow query warning
- Json column query and unmarshal to sub struct nested in data model
- GEO HASH 
- Nearby query by lng+lat+distance
- Built-in `db, protobuf, json` tag fetching  
- Read only column(s) specified by `sqlca:"readonly"` tag 
- Case...when syntax 
- Decimal, `sqlca.Decimal` instead of float64 for high precision calculation
- Where condition required when UPDATE/DELETE 
- [db2go](http://github.com/civet148/db2go) command line tool generate table schema output to .go or .proto file
- Custom sql.Scanner implement fetching 
- SSH tunnel
- Query results marshal to json string (QueryJson)
- Simply and developer friendly more than other ORM


# tests
[mysql schema](./examples/dbtest/test.sql)

[test cases](./examples/dbtest/main.go)


