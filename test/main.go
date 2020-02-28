package main

import "github.com/civet148/sqlca"

func main() {
	sqlca.NewEngine().Open(sqlca.AdapterSqlx_MySQL, "mysql://root:123456@tcp(127.0.0.1:3306)/enterprise?charset=utf8mb4")
}
