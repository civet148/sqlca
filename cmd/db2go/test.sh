#!/bin/sh
# options
# --without   exclude specify column(s)
# --readonly  read only column(s)

OUT_DIR=.
PACK_NAME=dataobject
go build && ./db2go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" \
--out $OUT_DIR --db "test" --table "users, classes" --suffix do --package $PACK_NAME --readonly "created_at, updated_at"

echo "go formatting..."
gofmt -w $OUT_DIR/$PACK_NAME
echo ok
