package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MysqlConf struct {
		Driver string
		Source string
	}
	RedisConf struct {
		Host     string
		Type     string `json:",default=node,options=node|cluster"`
		Pass     string `json:",optional"`
		Tls      bool   `json:",optional"`
		NonBlock bool   `json:",default=true"`
		// PingTimeout is the timeout for ping redis.
		PingTimeout time.Duration `json:",default=1s"`
	}
}
