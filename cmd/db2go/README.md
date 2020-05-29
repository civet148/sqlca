# db2go is a command to export database table structure to go or proto file 

## Usage

--url       connection url of database [required]

--db        databases to export, eg. "test" [required]

--table     tables to export, eg. "users, devices" [optional]

--out       output directory, default . [optional]

--package   export to a directory which used to be a package name in golang [optional]

--prefix    prefix of go file name [optional]

--suffix    suffix of go file name [optional]

--tag       customer golang struct member tag [optional] 

--without   ignore columns [optional]

--readonly  specify columns only for select [optional]

--proto     generate .proto file [optional]

--gogo-options specify gogo proto options [optional]

--one-file  all table schema integrate into one file named by database [optional]

```shell script
# db2go --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8" --db test --table users --package model --out $GOPATH/src/myproject
```
