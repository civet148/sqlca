package sqlca

import (
	"crypto/tls"
	"github.com/civet148/log"
	"time"
)

type RedisConfig struct {
	Address     string         // redis address, eg. "127.0.0.1:6379"
	Password    string         // redis password, default empty
	DB          int            // db index, default 0
	MaxActive   int            // max active connections
	MaxIdle     int            // max idle connections
	ConnTimeout *time.Duration // connection timeout

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration
	ClientName      string
	UseTLS          bool
	SkipVerify      bool
	TlsConfig       *tls.Config
}

type dialOption struct {
	Debug         bool         //enable debug mode
	Max           int          //max active connections
	Idle          int          //max idle connections
	SSH           *SSH         //ssh tunnel server config
	SnowFlake     *SnowFlake   //snowflake id config
	DisableOffset bool         //disable page offset for LIMIT (default page no is 1, if true then page no start from 0)
	DefaultLimit  int32        //limit default (0 means no limit)
	RedisConfig   *RedisConfig //redis config
}

type Option func(*dialOption)

var defaultDialOption = &dialOption{
	Max:  DefaultConnMax,
	Idle: DefaultConnIdle,
}

func parseDialOption(opts ...Option) *dialOption {
	for _, opt := range opts {
		opt(defaultDialOption)
	}
	log.Json(defaultDialOption)
	return defaultDialOption
}

func WithDebug() Option {
	return func(opt *dialOption) {
		opt.Debug = true
	}
}

func WithMaxConn(max int) Option {
	return func(opt *dialOption) {
		opt.Max = max
	}
}

func WithIdleConn(idle int) Option {
	return func(opt *dialOption) {
		opt.Idle = idle
	}
}

func WithDisableOffset() Option {
	return func(opt *dialOption) {
		opt.DisableOffset = true
	}
}

func WithDefaultLimit(limit int32) Option {
	return func(opt *dialOption) {
		opt.DefaultLimit = limit
	}
}

func WithSSH(ssh *SSH) Option {
	return func(opt *dialOption) {
		opt.SSH = ssh
	}
}

func WithSnowFlake(snowflake *SnowFlake) Option {
	return func(opt *dialOption) {
		opt.SnowFlake = snowflake
	}
}

func WithRedisConfig(rc *RedisConfig) Option {
	return func(opt *dialOption) {
		opt.RedisConfig = rc
	}
}
