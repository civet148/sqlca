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
	var results []map[string]string
	_, err := e.Model(&results).QueryMap("SELECT * FROM %v WHERE %v%v%v=%v%v%v",
		e.getTableName(),
		e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
		e.getSingleQuote(), e.getPkValue(), e.getSingleQuote())

	if err != nil {
		log.Errorf("%s", err)
		return
	}

	data, _ := json.Marshal(results) //marshal map to json string
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
		var results []map[string]string
		_, err := e.Model(&results).QueryMap("SELECT %v%v%v FROM %v WHERE %v%v%v=%v%v%v",
			e.getForwardQuote(), e.GetPkName(), e.getBackQuote(), e.getTableName(),
			e.getForwardQuote(), v.Name, e.getBackQuote(), e.getSingleQuote(), v.Value, e.getSingleQuote())
		if err != nil {
			log.Errorf("%s", err)
			return
		}

		if len(results) == 0 {
			continue
		}
		var pkValues []string
		for _, vv := range results {
			strCachePrimaryKey := e.makeCacheKey(e.GetPkName(), vv[e.GetPkName()])
			pkValues = append(pkValues, strCachePrimaryKey)
		}

		if len(pkValues) == 0 {
			continue
		}

		data, _ := json.Marshal(pkValues) //marshal []string to json string

		kv := &cacheKeyValue{
			Key: e.makeCacheKey(v.Name, v.Value), //index key in cache
			Value: cacheValue{
				ValueType:  ValueType_Index,
				TableName:  e.getTableName(),
				ColumnName: v.Name, //index column name in table
				CreateTime: strDateTime,
				ExpireSec:  e.expireTime,
				Data:       string(data), //pointer to primary key name in cache
			},
		}
		kvs = append(kvs, kv)
	}

	return
}

func (e *Engine) makeUpdateCache() (kvs []*cacheKeyValue) {

	if isNilOrFalse(e.getPkValue()) {
		log.Warnf("primary key's value is nil")
		return
	}
	kvs = append(kvs, e.makeCacheData())
	kvs = append(kvs, e.makeCacheIndexes()...)
	return
}

func (e *Engine) saveToCache(kvs ...*cacheKeyValue) (ok bool) {
	for _, v := range kvs {
		data, _ := json.Marshal(v.Value)
		if _, err := e.cache.Do("SETEX", v.Key, e.expireTime, string(data)); err != nil {
			log.Errorf("set key [%v] value [%v] error [%v]", v.Key, string(data), err.Error())
			return false
		} else {
			if e.isDebug() {
				e.getCacheValue(v.Key)
			}
		}
	}
	return true
}

func (e *Engine) upsertCache(lastInsertId int64) {

	if e.isCacheNil() && e.isDebug() {
		log.Debugf("cache instance is nil, can't update to cache")
		return
	}

	if !e.getUseCache() && e.isDebug() {
		log.Debugf("use cache is disabled, ignore it")
		return
	}

	if e.isPkInteger() {
		if lastInsertId == 0 {
			log.Errorf("got last insert id is 0")
			return
		}
		e.setPkValue(lastInsertId)
	}

	e.updateCache()
}

func (e *Engine) updateCache() (ok bool) {

	if e.isCacheNil() && e.isDebug() {
		log.Debugf("cache instance is nil, can't update to cache")
		return
	}

	if !e.getUseCache() && e.isDebug() {
		log.Debugf("use cache is disabled, ignore it")
		return
	}

	return e.saveToCache(e.makeUpdateCache()...)
}

func (e *Engine) queryCache() (count int64, ok bool) {

	if e.isCacheNil() && e.isDebug() {
		log.Warnf("cache instance is nil, can't update to cache")
		return
	}

	if !e.getUseCache() && e.isDebug() {
		log.Warnf("use cache is disabled, ignore it")
		return
	}

	if e.isPkValueNil() && len(e.getIndexes()) == 0 && e.isDebug() {
		log.Warnf("query condition primary key's value and index not set")
		return
	}

	if e.getOrderBy() != "" || e.getAscOrDesc() != "" || e.getGroupBy() != "" || e.getLimit() != "" || e.getOffset() != "" {
		if e.isDebug() {
			log.Warnf("query from cache can't use ORDER BY/GROUP BY/LIMIT/OFFSET key words")
		}
		return 0, false
	}

	if e.isPkValueNil() {
		return e.queryCacheByIndex()
	}
	return e.queryCacheById()
}

func (e *Engine) getCacheValue(strKey string) (kv *cacheKeyValue, ok bool) {
	kv = &cacheKeyValue{
		Key: strKey,
	}
	var err error
	var reply interface{}
	reply, err = e.cache.String(e.cache.Do("GET", strKey))
	if err = e.cache.Unmarshal(&kv.Value, reply, err); err != nil {
		log.Warnf("cache GET key [%v] error %v", strKey, err.Error())
		return kv, false
	}
	log.Debugf("cache GET key [%v] value %+v", strKey, kv.Value)
	return kv, true
}

func (e *Engine) queryCacheById() (count int64, ok bool) {

	strKey := e.makeCacheKey(e.GetPkName(), e.getPkValue())
	var kv *cacheKeyValue

	if kv, ok = e.getCacheValue(strKey); !ok {
		return 0, false
	}

	fetchers := e.makeCacheFetcher(kv)
	if len(fetchers) == 0 {
		log.Warnf("nothing found by id %+v value [%v] in cache", e.GetPkName(), e.getPkValue())
		return 0, false
	}

	return e.assignCacheToModel(fetchers)
}

func (e *Engine) queryCacheByIndex() (count int64, ok bool) {

	var fetchers []*Fetcher
	for _, v := range e.getIndexes() {

		var kv *cacheKeyValue
		strKey := e.makeCacheKey(v.Name, v.Value)

		if kv, ok = e.getCacheValue(strKey); !ok {
			log.Warnf("cache index GET [%v] value nil", strKey)
			return
		}

		var cacheIdx cacheIndex
		if err := json.Unmarshal([]byte(kv.Value.Data), &cacheIdx.Keys); err != nil {
			log.Errorf("unmarshal [%v] to cache index error [%v]", kv.Value.Data, err.Error())
			return
		}

		for _, vv := range cacheIdx.Keys {
			//log.Debugf("cache GET [%v] ready", vv)
			if kv, ok = e.getCacheValue(vv); !ok {
				log.Warnf("cache index GET [%v] value nil", vv)
				return 0, false
			}
			f := e.makeCacheFetcher(kv)
			fetchers = append(fetchers, f...)
		}
	}

	if len(fetchers) == 0 {
		log.Warnf("nothing found by index %+v in cache", e.getIndexes())
		return 0, false
	}

	return e.assignCacheToModel(fetchers)
}

func (e *Engine) assignCacheToModel(fetchers []*Fetcher) (count int64, ok bool) {

	var err error
	//log.Debugf("assign cache data to fetchers")
	if e.getModelType() == ModelType_BaseType {
		if count, err = e.fetchCache(fetchers, e.model.([]interface{})...); err != nil {
			log.Errorf("fetchRow error [%v]", err.Error())
			return 0, false
		}
	} else {
		if count, err = e.fetchCache(fetchers, e.model); err != nil {
			log.Errorf("fetchRow error [%v]", err.Error())
			return 0, false
		}
	}
	//log.Debugf("model [%+v] count [%v] ok [%v]", e.model, count, ok)
	return count, true
}

func (e *Engine) makeCacheFetcher(kvs ...*cacheKeyValue) (fetchers []*Fetcher) {

	for _, v := range kvs {
		var mapValues []map[string]string
		if err := json.Unmarshal([]byte(v.Value.Data), &mapValues); err != nil {
			log.Errorf("cache value [%v] unmarshal to map[string]string error [%+v]", v.Value, err.Error())
			continue
		}

		for _, vv := range mapValues {
			count := len(vv)
			fetcher := &Fetcher{
				count:     count,
				cols:      nil,
				types:     nil,
				arrValues: mapToBytesSlice(vv),
				mapValues: vv,
				arrIndex:  0,
			}
			fetchers = append(fetchers, fetcher)
		}
	}

	return
}
