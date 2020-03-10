@echo off
go run main.go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" ^
--out "%GOPATH%/src/nebula.chat/enterprise/bot/dal" ^
--db "kefu_system" --table "assignments"  --prefix kefu --suffix do --package dataobject