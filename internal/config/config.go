package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MysqlConf MysqlConf
	RedisConf RedisConf
}

type MysqlConf struct {
	Driver string
	Source string
}

type RedisConf struct {
	Addr        string
	Type        string `json:",default=node,options=node|cluster"`
	DB          int
	Requirepass string `json:",optional"`
	Tls         bool   `json:",optional"`
	NonBlock    bool   `json:",default=true"`
	// PingTimeout is the timeout for ping redis.
	PingTimeout time.Duration `json:",default=1s"`
}
