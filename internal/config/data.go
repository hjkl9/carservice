package config

import "time"

type Data struct {
	Mysql MysqlConf
	Redis RedisConf
}

type MysqlConf struct {
	Driver string
	Source string
}

type RedisConf struct {
	Host     string
	Type     string `json:",default=node,options=node|cluster"`
	Pass     string `json:",optional"`
	Tls      bool   `json:",optional"`
	NonBlock bool   `json:",default=true"`
	// PingTimeout is the timeout for ping redis.
	PingTimeout time.Duration `json:",default=1s"`
}
