package sqlca

import (
	"github.com/astaxie/beego/cache"
	"github.com/civet148/gotools/log"
)

type ValueType int

const (
	CACHE_INDEX_DEEP = 1 // index deep in cache
)

const (
	ValueType_Data  ValueType = 1 // data of table
	ValueType_Index ValueType = 2 // index of data
)

func (v ValueType) GoString() string {
	return v.String()
}

func (v ValueType) String() string {
	switch v {
	case ValueType_Data:
		return "ValueType_Data"
	case ValueType_Index:
		return "ValueType_Index"
	}
	return "ValueType_Unknown"
}

type CacheValue struct {
	ValueType ValueType `json:"value_type"` // cache value type
	TableName string    `json:"table_name"` // table name
	CreatedAt string    `json:"created_at"` // cache data create time
	ExpiredAt string    `json:"expired_at"` // cache data expire time
	Data      string    `json:"data"`       // index or data json in redis/memcached...
}

type CacheIndex struct {
	Keys []string `json:"keys"`
}

// [redis]    newCache("redis",    `{"conn":"127.0.0.1:6379"}`)
// [memcache] newCache("memcache", `{"conn":"127.0.0.1:11211"}`)
// [memory]   newCache("memory",   `{"interval":60}`)
func newCache(name, config string) (c cache.Cache, err error) {
	c, err = cache.NewCache(name, config)
	if err != nil {
		log.Errorf("%v", err.Error())
		return
	}
	return
}

func (c *CacheValue) getCacheDataKey() (strKey string) {

	return
}

func (c *CacheValue) getCacheIndexKey() (strKey string) {

	return
}
