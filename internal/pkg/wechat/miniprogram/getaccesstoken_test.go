package miniprogram

import (
	"carservice/internal/config"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	cfg := config.WechatConf{
		// ! Should recover to empty.
		MiniProgram: config.MiniProgram{
			AppId:  "",
			Secret: "",
		},
	}
	mp := NewWechatProvider(cfg)
	mp.GetAccessToken()
}
