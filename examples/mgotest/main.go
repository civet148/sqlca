package main

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3"
)

func main() {
	//mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	url := "mongodb://admin:123456@127.0.0.1:27017/test?authSource=admin"

	e, err := sqlca.NewEngine(url)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	e.Debug(true) //debug on

}
