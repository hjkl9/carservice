package config

import (
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	LogConf logx.LogConf

	MysqlConf MysqlConf
	RedisConf RedisConf

	JwtConf    JwtConf
	SmsConf    SmsConf
	WechatConf WechatConf
}

type LogConf struct {
	ServiceName string `json:",optional"`
	Mode        string `json:",default=console,options=[console,file,volume]"`
	Encoding    string `json:",default=json,options=[json,plain]"`
	TimeFormat  string `json:",optional"`
	Path        string `json:",default=logs"`
	Level       string `json:",default=info,options=[debug,info,error,severe]"`
	// MaxContentLength    uint32 `json:",optional"`
	Compress            bool   `json:",optional"`
	Stat                bool   `json:",default=true"`
	KeepDays            int    `json:",optional"`
	StackCooldownMillis int    `json:",default=100"`
	MaxBackups          int    `json:",default=0"`
	MaxSize             int    `json:",default=0"`
	Rotation            string `json:",default=daily,options=[daily,size]"`
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
	MiniProgram MiniProgram
}

type MiniProgram struct {
	AppId  string
	Secret string
}
