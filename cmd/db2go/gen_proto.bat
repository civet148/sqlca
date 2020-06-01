@echo off
go build -ldflags "-s -w"

set OUT_DIR=.
set PACK_NAME=proto
set GOGO_OPTIONS="(gogoproto.marshaler_all)=true,(gogoproto.sizer_all)=true,(gogoproto.unmarshaler_all)=true,(gogoproto.gostring_all)=true"

go run main.go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" --disable-decimal --proto --gogo-options %GOGO_OPTIONS% ^
--out %OUT_DIR% --db "test" --table "users, classes" --suffix do --package %PACK_NAME% --one-file

echo generate protobuf file ok
pause