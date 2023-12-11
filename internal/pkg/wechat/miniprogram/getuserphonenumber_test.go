package miniprogram

import (
	"carservice/internal/config"
	"testing"
)

func TestGetUserPhoneNumber(t *testing.T) {
	cfg := config.WechatConf{
		// ! Should recover to empty.
		MiniProgram: config.MiniProgram{
			AppId:  "",
			Secret: "",
		},
	}
	provider := NewWechatProvider(cfg)
	provider.GetUserPhoneNumber("faker code.", "")
}
