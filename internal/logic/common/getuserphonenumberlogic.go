package common

import (
	"context"
	"net/http"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/wechat"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPhoneNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPhoneNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPhoneNumberLogic {
	return &GetUserPhoneNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPhoneNumberLogic) GetUserPhoneNumber(req *types.GetUserPhoneNumberReq) (resp *types.GetUserPhoneNumberRep, err error) {
	// get wechat util.
	wx := wechat.NewWechat(l.svcCtx.Config.WechatConf)
	// get mini program from wechat util.
	mp := wx.MiniProgram()
	// get access token.
	accessToken, err := mp.GetAccessToken()
	if err != nil {
		return nil, errcode.New(http.StatusOK, "", err.Error())
	}
	// get phone number by code and access token.
	phoneNumber, err := mp.GetUserPhoneNumber(accessToken, req.Code)
	if err != nil {
		return nil, errcode.New(http.StatusOK, "", err.Error())
	}
	// response
	return &types.GetUserPhoneNumberRep{
		PhoneNumber: phoneNumber,
	}, nil
}
