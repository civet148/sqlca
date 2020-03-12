# db2go is a command to export database table structure to go types 

## Usage

parameters
--url       connection url of database

--db        databases to export

--table     tables to export

--out       output directory, default .

--package   export to a directory which used to be a package name in golang

--prefix    prefix of go file name

--suffix    suffix of go file name

--tag       customer golang struct member tag 

```shell script
# db2go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" --db test --table users --package model --out $GOPATH/src/myproject
```
