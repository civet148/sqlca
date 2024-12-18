@echo off


rem 数据模型(models)和数据库操作对象(dao)文件输出基础目录
set OUT_DIR=.
rem 数据模型包名(数据模型文件目录名)
set PACK_NAME="models"
rem 指定某表的某字段为指定类型,多个表字段以英文逗号分隔（例如：user.create_time=time.Time表示指定user表create_time字段为time.Time类型; 如果不指定表名则所有表的create_time字段均为time.Time类型；支持第三方包类型，例如：user.weight=github.com/shopspring/decimal.Decimal）
set SPEC_TYPES="inventory_data.product_extra=ProductExtraData"
rem 指定其他orm的标签和值(以空格分隔)
set COMMON_TAGS="id=gorm:\"primarykey\" create_time=gorm:\"autoCreateTime\" update_time=gorm:\"autoUpdateTime\""
set DEPLOY_SQL="test.sql"

rem 判断本地系统是否已安装db2go工具，没有则进行安装
echo "searching db2go.exe ..."
echo "--------------------------------------------"
where db2go.exe
echo "--------------------------------------------"

IF "%errorlevel%" == "0" (
    echo db2go already installed.
) ELSE (
    echo db2go not found in system %%PATH%%, installing...
    go install github.com/civet148/db2go@latest
    If "%errorlevel%" == "0" (
        echo db2go install successfully.
    ) ELSE (
        rem 安装失败: Linux/Mac请安装gcc工具链，Windows系统可以安装msys64进行源码编译或通过链接直接下载二进制(最新版本v2.13 https://github.com/civet148/release/tree/master/db2go/v2)
        echo ERROR: Linux/Mac please install gcc tool-chain and windows download from https://github.com/civet148/release/tree/master/db2go/v2 (latest version is v2.13)
    )
)

rem ---------------------- 导出数据库表结构-------------------------
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
db2go   --url %DSN_URL% --out %OUT_DIR% --spec-type %SPEC_TYPES% --package %PACK_NAME%  --common-tags %COMMON_TAGS% --export %DEPLOY_SQL%
gofmt -w %OUT_DIR%/%PACK_NAME%

pause


