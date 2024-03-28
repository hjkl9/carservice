package userorder

import (
	"context"

	"carservice/internal/pkg/wechat/payment"
	"carservice/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentOrderCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentOrderCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentOrderCallbackLogic {
	return &PaymentOrderCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentOrderCallbackLogic) PaymentOrderCallback(req *payment.DecryptedResource) error {

	return nil
}
