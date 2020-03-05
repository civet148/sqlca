package sqlca

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/redigogo"
	_ "github.com/civet148/redigogo/alone"
	_ "github.com/civet148/redigogo/cluster"
	"time"
)

type valueType int

const (
	CACHE_INDEX_DEEP  = 1           // index deep in cache
	CACHE_REPLICATE   = "replicate" //replicate host [ip:port,...]
	CACHE_DB_INDEX    = "db"
	CAHCE_SQLX_PREFIX = "sqlx:cache"
)

const (
	ValueType_Data  valueType = 1 // data of table
	ValueType_Index valueType = 2 // index of data
)

func (v valueType) GoString() string {
	return v.String()
}

func (v valueType) String() string {
	switch v {
	case ValueType_Data:
		return "ValueType_Data"
	case ValueType_Index:
		return "ValueType_Index"
	}
	return "ValueType_Unknown"
}

type cacheValue struct {
	ValueType  valueType `json:"value_type"`  // cache value type
	TableName  string    `json:"table_name"`  // table name
	ColumnName string    `json:"column_name"` // table column name
	CreatedAt  string    `json:"created_at"`  // cache data create time
	ExpiredAt  string    `json:"expired_at"`  // cache data expire time
	Data       string    `json:"data"`        // index or data json in redis/memcached...
}

type cacheIndex struct {
	Keys []string `json:"keys"`
}

type cacheKeyValue struct {
	Key   string     `json:"key"`
	Value cacheValue `json:"value"`
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

//
//func (c *CacheValue) makeIndexKey() string {
//	return
//}

func (e *Engine) makeCacheKey(name string, value interface{}) string {
	return fmt.Sprintf("%v:%v:%v:%v", CAHCE_SQLX_PREFIX, e.getTableName(), name, value)
}

func (e *Engine) marshalModel() (s string) {

	//data, _ := json.Marshal(e.model)
	return //string(data)
}

func (e *Engine) unmarshalModel() (s string) {

	return
}

func (e *Engine) getDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (e *Engine) queryCacheData(strCondition string) (res []map[string]string, err error) {

	strQuery := fmt.Sprintf("SELECT * FROM %v WHERE %v", e.getTableName(), strCondition)
	log.Debugf("query sql [%v]", strQuery)
	return e.QueryMap(strQuery)
}

func (e *Engine) makeCacheData() (kv *cacheKeyValue) {

	//get current datetime
	strDateTime := e.getDateTime()
	//make cache key and data
	strDataKey := e.makeCacheKey(e.GetPkName(), e.getPkValue())

	return &cacheKeyValue{
		Key: strDataKey,
		Value: cacheValue{
			ValueType:  ValueType_Data,
			TableName:  e.getTableName(),
			ColumnName: e.GetPkName(),
			CreatedAt:  strDateTime,
			ExpiredAt:  strDateTime,
			Data:       e.marshalModel(),
		},
	}
}

func (e *Engine) makeCacheIndexes() (kvs []*cacheKeyValue) {

	//get current datetime
	strDateTime := e.getDateTime()
	//make cache key and data
	strDataKey := e.makeCacheKey(e.GetPkName(), e.getPkValue())

	for _, v := range e.getIndexes() {

		kv := &cacheKeyValue{
			Key: e.makeCacheKey(v.name, v.value), //index key in cache
			Value: cacheValue{
				ValueType:  ValueType_Index,
				TableName:  e.getTableName(),
				ColumnName: v.name, //index column name in table
				CreatedAt:  strDateTime,
				ExpiredAt:  strDateTime,
				Data:       strDataKey, //pointer to primary key name in cache
			},
		}
		kvs = append(kvs, kv)
	}

	return
}

func (e *Engine) makeCache() (kvs []*cacheKeyValue) {

	assert(e.pkValue, "primary key [%v] value is nil, please call Id() method", e.GetPkName())

	kvs = append(kvs, e.makeCacheData())
	kvs = append(kvs, e.makeCacheIndexes()...)
	log.Debugf("makeCache kvs [%+v]", kvs)
	return
}
