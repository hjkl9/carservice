package wechat

import (
	"carservice/internal/config"
	"testing"
)

func TestNewWechat(t *testing.T) {
	cfg := config.WechatConf{
		MiniProgram: config.MiniProgram{
			AppId:  "",
			Secret: "",
		},
	}
	wechat := NewWechat(cfg)
	wechat.MiniProgram().GetUserPhoneNumber("accessToken", "code")
}
