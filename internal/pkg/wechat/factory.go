package wechat

import (
	"carservice/internal/config"
	"carservice/internal/pkg/wechat/miniprogram"
)

type WechatProvider interface {
	MiniProgram() miniprogram.MiniProgram
}

func NewWechatProvider(config config.WechatConf) WechatProvider {
	return NewWechat(config)
}

type wechat struct {
	config config.WechatConf
}

func NewWechat(config config.WechatConf) *wechat {
	return &wechat{
		config,
	}
}

func (w *wechat) MiniProgram() miniprogram.MiniProgram {
	return miniprogram.NewWechatProvider(w.config)
}
