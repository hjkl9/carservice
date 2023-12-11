package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MysqlConf MysqlConf
	RedisConf RedisConf

	JwtConf    JwtConf
	SmsConf    SmsConf
	WechatConf WechatConf
}

type MysqlConf struct {
	Driver string
	Source string
}

type RedisConf struct {
	Addr        string
	DB          int
	Requirepass string `json:",optional"`
	// PingTimeout is the timeout for ping redis.
	PingTimeout time.Duration `json:",default=1s"`
}

type JwtConf struct {
	AccessSecret string
	AccessExipre int `json:",optional"`
}

type SmsConf struct {
	SecretId   string
	SecretKey  string
	SdkAppId   string
	SignName   string
	TemplateId string
}

type WechatConf struct {
	MiniProgram struct {
		AppId  string
		Secret string
	}
}
