#!/bin/sh
# options
# --without   exclude specify column(s)
# --readonly  read only column(s)
go build && ./db2go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" \
--out . --db "test" --table "users, classes" --suffix do --package dataobject --readonly "created_at, updated_at"