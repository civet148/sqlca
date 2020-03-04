package sqlca

import (
	"encoding/json"
	redis "github.com/gitstliu/go-redis-cluster"
	"time"
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

type CacheConfig struct {
	Password       string   `json:"password"`
	Index          int      `json:"db_index"`
	MasterHost     string   `json:"master_host"`
	ReplicateHosts []string `json:"replicate_hosts"`
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

type Cache struct {
	c *redis.Cluster
}

func newCache(strScheme string, strConfig string) (cache *Cache, err error) {

	var config CacheConfig
	if err = json.Unmarshal([]byte(strConfig), &config); err != nil {
		assert(false, "cache config [%v] illegal", strConfig)
	}

	var StartNodes []string
	StartNodes = append(StartNodes, config.MasterHost)
	StartNodes = append(StartNodes, config.ReplicateHosts...)

	var c *redis.Cluster
	c, err = redis.NewCluster(&redis.Options{
		StartNodes:   StartNodes,
		ConnTimeout:  500 * time.Millisecond,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		KeepAlive:    16,
		AliveTime:    60 * time.Second,
	})

	if err != nil {
		assert(false, "redis cluster %v connect error [%v]", StartNodes, err.Error())
	}
	return &Cache{c: c}, nil
}

func (c *CacheValue) getCacheDataKey() (strKey string) {

	return
}

func (c *CacheValue) getCacheIndexKey() (strKey string) {

	return
}
