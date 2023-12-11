package wechat

import "carservice/internal/pkg/wechat/miniprogram"

type WechatProvider interface {
	MiniProgram() miniprogram.MiniProgramProvider
}
