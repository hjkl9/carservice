package userorder

import (
	"context"

	"carservice/internal/datatypes/carreplacement"
	enum_userorder "carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/pkg/wechat/payment"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
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

func (l *PaymentOrderLogic) PaymentOrder(req *types.PaymentOrderReq) (*types.PaymentOrderRep, error) {
	userId := jwt.GetUserId(l.ctx)
	orderId := req.Id
	replacementIds := req.CarReplacements
	// ? 处理和计算客户端的配件列表价格
	l.filterAndCalcAmount()
	// ! 删除临时配件
	// replacementIds = []int{5, 6, 14}

	// 检查订单是否存在
	hasOrder, err := l.svcCtx.Repo.UserOrderRelated().GetIfOrderExistsById(l.ctx, userId, uint(orderId))
	if err != nil {
		return nil, errcode.DatabaseGetErr
	}
	if !hasOrder {
		return nil, errcode.OrderNotFoundErr
	}

	// 获取订单
	order, err := l.svcCtx.Repo.
		UserOrderRelated().
		GetOrderById(l.ctx, userId, uint(orderId))
	if err != nil {
		logc.Errorf(l.ctx, "获取订单时发生错误, err: %s\n", err.Error())
		return nil, errcode.DatabaseGetErr
	}

	// 获取用户 open_id
	openId, err := l.svcCtx.Repo.UserRelated().GetOpenIdByUserId(l.ctx, userId)
	if err != nil {
		return nil, errcode.DatabaseGetErr
	}

	if uint8(order.OrderStatus) != enum_userorder.AwaitingPayment {
		return nil, errcode.OrderOprErr.SetMessage("无法支付非待支付的订单")
	}

	// 获取车型价格和匹配
	officialPriceDown, officialPriceUp, err := l.svcCtx.Repo.
		CarBrandSeriesRepo().
		GetOfficialPrice(l.ctx, order.CarBrandSeriesId)
	if err != nil {
		logc.Errorf(l.ctx, "查询车型官方报价发生错误, err: %s\n", err.Error())
		return nil, errcode.DatabaseGetErr.SetMessage("订单官方报价查询发生错误")
	}
	// todo: 处理未报价的车型
	if officialPriceDown == officialPriceUp {
	}

	// True: 高端
	// False: 低端
	var gradeFunc func() bool = func() bool {
		return l.svcCtx.Repo.
			CarBrandSeriesRepo().
			CheckGradeByCarSeries(officialPriceDown, officialPriceUp)
	}

	// 获取配件列表
	replacements, err := l.svcCtx.Repo.
		CarReplacementRepoRelated().
		GetEstPriceListByIdSet(l.ctx, gradeFunc, replacementIds)
	if err != nil {
		logc.Errorf(l.ctx, "获取配件列表发生错误, err: %s\n", err.Error())
		return nil, errcode.DatabaseGetErr
	}

	// 匹配和计算配件价格
	totalEstF32Price, totalEstU64Price := l.calcAmount(replacements)
	var totalAmount struct {
		estF32Price float64
		estU64Price uint64
	}
	totalAmount.estF32Price = totalEstF32Price
	totalAmount.estU64Price = totalEstU64Price

	// 准备预支付数据
	payload := payment.PaymentPayload{
		Description: "TODO",
		OutTradeNo:  order.OrderNumber,
		Attach:      "TODO",
		NotifyUrl:   l.svcCtx.Config.AppUrl + "/v1/userOrder/pay/callback",
		Amount:      int64(totalAmount.estU64Price), // 一分钱
		OpenId:      openId,
	}

	_ = payload

	return &types.PaymentOrderRep{
		GoSdkVersion: "",
		PrepayId:     "",
		TimeStamp:    "",
		NonceStr:     "",
		Package:      "",
		SignType:     "",
		PaySign:      "",
	}, nil

	// 准备支付配置
	ourConf := l.svcCtx.Config.WechatPayMerchantConf
	cfg := payment.PaymentConfig{
		MchId:               ourConf.MchId,
		MchCertSerialNumber: ourConf.MchCertSerialNumber,
		MchApiV3Key:         ourConf.MchApiV3Key,
		Appid:               ourConf.AppId,
		PrivateKeyPath:      ourConf.PvtKeyPath,
	}
	paymentResp, err := payment.PrepayOrder(cfg, payload)
	if err != nil {
		return nil, errcode.OrderPrepayErr.WithDetails(err.Error())
	}

	return &types.PaymentOrderRep{
		GoSdkVersion: "v0.2.18",
		PrepayId:     *paymentResp.PrepayId,
		TimeStamp:    *paymentResp.TimeStamp,
		NonceStr:     *paymentResp.NonceStr,
		Package:      *paymentResp.Package,
		SignType:     *paymentResp.SignType,
		PaySign:      *paymentResp.PaySign,
	}, nil
}

func (l *PaymentOrderLogic) calcAmount(
	replacements []carreplacement.Replacement,
) (float64, uint64) {
	var f float32 = 0.00
	var u uint64 = 0
	for _, replacement := range replacements {
		f = f + replacement.EstF32Price
		u = u + replacement.EstU64Price
	}
	return float64(f), u
}

func (l *PaymentOrderLogic) filterAndCalcAmount() {

}
