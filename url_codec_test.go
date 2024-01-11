package sqlca

import (
	"github.com/civet148/log"
	"testing"
)

func TestMySqlRawDSN(t *testing.T) {
	dsn, err := Url2MySql("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")
	if err != nil {
		t.Error(err)
		return
	}
	log.Infof(dsn)
}
