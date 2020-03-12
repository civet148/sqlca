# db2go is a command to export database table structure to go types 

## Usage

parameters
--url       connection url of database [required]

--db        databases to export, eg. "db1,db2" [required]

--table     tables to export, eg. "users, devices" [optional]

--out       output directory, default . [required]

--package   export to a directory which used to be a package name in golang [optional]

--prefix    prefix of go file name [optional]

--suffix    suffix of go file name [optional]

--tag       customer golang struct member tag [optional] 

```shell script
# db2go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" --db test --table users --package model --out $GOPATH/src/myproject
```
