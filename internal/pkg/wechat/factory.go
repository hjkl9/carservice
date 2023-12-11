package wechat

import (
	"carservice/internal/config"
	"carservice/internal/pkg/wechat/miniprogram"
)

type WechatProvider interface {
	MiniProgram() miniprogram.MiniProgram
}

type Wechat struct {
	config config.WechatConf
}

func NewWechat(config config.WechatConf) *Wechat {
	return &Wechat{
		config,
	}
}

func (w *Wechat) MiniProgram() miniprogram.MiniProgram {
	return miniprogram.NewWechatProvider(w.config)
}
