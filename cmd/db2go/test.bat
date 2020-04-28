@echo off
go build
db2go.exe --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" ^
--out . --db "test" --table "users, classes" --suffix do --package dataobject --readonly "created_at, updated_at"
pause