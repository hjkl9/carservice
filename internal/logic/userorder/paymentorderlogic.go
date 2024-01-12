package userorder

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentOrderLogic {
	return &PaymentOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentOrderLogic) PaymentOrder(req *types.PaymentOrderReq) (*types.PaymentOrderReq, error) {
	// todo: add your logic here and delete this line

	return nil, nil
}
