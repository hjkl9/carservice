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

	JwtConf    JwtConf    // 认证配置
	SmsConf    SmsConf    // 短信相关配置
	WechatConf WechatConf // 微信相关配置
	AMapConf   AMapConf   // 高德地图配置
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
	AccessExpire int `json:",optional"`
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

type AMapConf struct {
	Key string
}
