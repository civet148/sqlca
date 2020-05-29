@echo off

set SRC_DIR=.\
set DST_DIR=.\

echo generating...

protoc -I=%SRC_DIR%  --proto_path=%GOPATH%\src  --gogo_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com\gogo\protobuf\types,:%DST_DIR%  %SRC_DIR%\*.proto

echo generate over

pause
