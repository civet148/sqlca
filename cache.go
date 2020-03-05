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
		return "Data"
	case ValueType_Index:
		return "Index"
	}
	return "Unknown"
}

type cacheValue struct {
	ValueType  valueType `json:"value_type"`  // cache value type
	TableName  string    `json:"table_name"`  // table name
	ColumnName string    `json:"column_name"` // table column name
	CreateTime string    `json:"create_time"` // cache data create time
	ExpireSec  int       `json:"expire_sec"`  // cache data expire time
	Data       string    `json:"data"`        // index or data json in redis
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

func (e *Engine) makeCacheData() (kv *cacheKeyValue) {

	//get current datetime
	strDateTime := e.getDateTime()
	//make cache key and data
	strPrimaryCacheKey := e.makeCacheKey(e.GetPkName(), e.getPkValue())
	m, _ := e.QueryMap("SELECT * FROM %v WHERE %v%v%v=%v%v%v",
		e.getTableName(),
		e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
		e.getSingleQuote(), e.getPkValue(), e.getSingleQuote())
	data, _ := json.Marshal(m) //marshal map to json string
	return &cacheKeyValue{
		Key: strPrimaryCacheKey,
		Value: cacheValue{
			ValueType:  ValueType_Data,
			TableName:  e.getTableName(),
			ColumnName: e.GetPkName(),
			CreateTime: strDateTime,
			ExpireSec:  e.expireTime,
			Data:       string(data),
		},
	}
}

func (e *Engine) makeCacheIndexes() (kvs []*cacheKeyValue) {

	//get current datetime
	strDateTime := e.getDateTime()
	//make cache key and data
	for _, v := range e.getIndexes() {
		//eg. "select `id` from users where `phone`='8615439905001'"
		//`id` is primary key and `phone` is an index (maybe exist multiple records)
		m, _ := e.QueryMap("SELECT %v%v%v FROM %v WHERE %v%v%v=%v%v%v",
			e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
			e.getTableName(),
			e.getForwardQuote(), v.name, e.getBackQuote(),
			e.getSingleQuote(), v.value, e.getSingleQuote())

		var pkValues []string
		for _, vv := range m {
			strCachePrimaryKey := e.makeCacheKey(e.GetPkName(), vv[e.GetPkName()])
			pkValues = append(pkValues, strCachePrimaryKey)
		}
		data, _ := json.Marshal(pkValues) //marshal []string to json string

		kv := &cacheKeyValue{
			Key: e.makeCacheKey(v.name, v.value), //index key in cache
			Value: cacheValue{
				ValueType:  ValueType_Index,
				TableName:  e.getTableName(),
				ColumnName: v.name, //index column name in table
				CreateTime: strDateTime,
				ExpireSec:  e.expireTime,
				Data:       string(data), //pointer to primary key name in cache
			},
		}
		kvs = append(kvs, kv)
	}

	return
}

func (e *Engine) makeCache() (kvs []*cacheKeyValue) {

	if isNilOrFalse(e.getPkValue()) {
		e.setPkValue(e.getModelValue(e.GetPkName()))
	}

	assert(e.getPkValue(), "primary key [%v] value is nil, please call Id() method", e.GetPkName())

	kvs = append(kvs, e.makeCacheData())
	kvs = append(kvs, e.makeCacheIndexes()...)
	return
}

func (e *Engine) updateCache() {

	if e.isCacheNil() && e.isDebug() {
		log.Debugf("cache instance is nil, can't update to cache")
		return
	}

	if !e.getUseCache() && e.isDebug() {
		log.Debugf("use cache is disabled, ignore it")
		return
	}

	kvs := e.makeCache()
	for _, v := range kvs {
		data, _ := json.Marshal(v.Value)
		if _, err := e.cache.Do("SETEX", v.Key, e.expireTime, string(data)); err != nil {
			log.Errorf("set key [%v] value [%v] error [%v]", v.Key, string(data), err.Error())
		} else {
			//test read from cache
			kv := cacheKeyValue{
				Key: v.Key,
			}
			var reply interface{}
			reply, err = e.cache.String(e.cache.Do("GET", v.Key))
			if err = e.cache.Unmarshal(&kv.Value, reply, err); err != nil {
				log.Errorf("cache GET key [%v] error %v", v.Key, err.Error())
			} else {
				log.Debugf("cache GET key [%v] value %+v", v.Key, kv.Value)
			}
		}
	}
}
