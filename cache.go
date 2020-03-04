package sqlca

import (
	"encoding/json"
	"github.com/civet148/redigogo"
	_ "github.com/civet148/redigogo/alone"
	_ "github.com/civet148/redigogo/cluster"
)

type ValueType int

const (
	CACHE_INDEX_DEEP = 1           // index deep in cache
	CACHE_REPLICATE  = "replicate" //replicate host [ip:port,...]
	CACHE_DB_INDEX   = "db"
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

func newCache(strScheme string, strConfig string) (cache redigogo.Cache, err error) {

	var config redigogo.Config
	if err = json.Unmarshal([]byte(strConfig), &config); err != nil {
		assert(false, "cache config [%v] illegal", strConfig)
	}

	cache = redigogo.NewCache(&config)
	if err != nil {
		assert(false, "redis cache with config [%+v] create error [%v]", config, err.Error())
	}
	return
}

func (c *CacheValue) getCacheDataKey() (strKey string) {

	return
}

func (c *CacheValue) getCacheIndexKey() (strKey string) {

	return
}
