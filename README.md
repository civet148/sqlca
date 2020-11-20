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
- Batch insert
- Query/Insert/Update
- Upsert by customization
- Transactions wrapper
- Slow query warn
- Json column query and unmarshal to sub struct nested in data model object
- db2go commander line tool generate table schema output to .go or .proto file
- GEO HASH 
- Nearby query by lng+lat+distance
- Multiple databases (MySQL/Postgres/MS-SQLSERVER), read/write splitting  
- Read only column specified by `sqlca:"readonly"` tag 
- Simply and developer friendly more than other ORM


# test cases
table schema test.sql

see test/main.go


