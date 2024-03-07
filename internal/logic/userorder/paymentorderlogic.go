package userorder

import (
	"context"
	"fmt"

	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/conv"
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

type Replacement struct {
}

// todo: 统一封装
type OfficialPrice struct {
	PriceUp   float64 `db:"priceUp"`
	PriceDown float64 `db:"priceDown"`
}

func (l *PaymentOrderLogic) PaymentOrder(req *types.PaymentOrderReq) (*types.PaymentOrderRep, error) {
	userId := jwt.GetUserId(l.ctx)
	orderId := req.Id
	replacementIds := req.CarReplacements
	// ! 删除临时配件
	replacementIds = []int{5, 6, 14}

	hasOrder, err := l.svcCtx.Repo.UserOrderRelated().GetIfOrderExistsById(l.ctx, userId, uint(orderId))
	if err != nil {
		return nil, errcode.DatabaseGetErr
	}
	if !hasOrder {
		return nil, errcode.OrderNotFoundErr
	}

	var order struct {
		CarBrandSeriesId int64 `db:"carBrandSeriesId"`
		OrderStatus      uint8 `db:"orderStatus"`
	}

	// 获取订单
	query := "SELECT `car_brand_series_id` AS `carBrandSeriesId`, `order_status` AS `orderStatus` FROM `user_orders` WHERE `id` = ? LIMIT 1;"
	stmtx, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return nil, errcode.DatabasePrepareErr
	}
	if err = stmtx.GetContext(l.ctx, &order, orderId); err != nil {
		fmt.Println(err.Error())
		return nil, errcode.DatabaseGetErr
	}

	if order.OrderStatus != userorder.AwaitingPayment {
		return nil, errcode.OrderOprErr.SetMessage("无法支付非待支付的订单")
	}

	// 获取车型价格和匹配
	var officialPrice OfficialPrice
	query = "SELECT `official_price_up` AS `priceUp`, `official_price_down` AS `priceDown` FROM `car_brand_series` WHERE `series_id` = ? LIMIT 1;"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "查询官方售价[预处理时]发生错误", err)
		return nil, errcode.DatabasePrepareErr
	}
	if err = stmt.GetContext(l.ctx, &officialPrice, order.CarBrandSeriesId); err != nil {
		logc.Error(l.ctx, "查询官方售价[获取数据时]发生错误", err)
		return nil, errcode.DatabaseGetErr
	}

	// True: 高端
	// False: 低端
	var grade bool = func() bool {
		const P float64 = 30.00
		if officialPrice.PriceDown > P || officialPrice.PriceUp > P {
			return true
		}
		// 其他条件或规则
		// todo: 处理过滤规则
		return false
	}()
	var addQuery = func() string {
		if grade {
			return "`hm_est_f32_price` AS `estF32Price`, `hm_est_u64_price` AS `estU64Price`"
		} else {
			return "`lm_est_f32_price` AS `estF32Price`, `lm_est_u64_price` AS `estU64Price`"
		}
	}()
	query = "SELECT %s FROM `car_replacements` WHERE `id` IN (?);"
	query = fmt.Sprintf(query, addQuery)
	var replacements []struct {
		EstF32Price float64 `db:"estF32Price"`
		EstU64Price uint64  `db:"estU64Price"`
	}
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "查询配件列表[预处理时]发生错误", err)
		return nil, errcode.DatabasePrepareErr
	}
	replacementsStr := conv.ToStringWithSep_int(',', replacementIds...)
	err = stmt.SelectContext(l.ctx, &replacements, replacementsStr)
	if err != nil {
		logc.Error(l.ctx, "查询配件列表[获取数据时]发生错误", err)
		return nil, errcode.DatabaseGetErr
	}
	fmt.Println(replacements)
	return &types.PaymentOrderRep{
		Comment:  "*paymentResp.NonceStr",
		PrepayId: "*paymentResp.PrepayId",
	}, nil

	// 获取配件列表

	// 匹配和计算配件价格

	// 准备预支付数据
	payload := payment.PaymentPayload{
		Description: "",
		OutTradeNo:  "",
		Attach:      "",
		NotifyUrl:   "",
		Amount:      1, // 一分钱
		OpenId:      "",
	}

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
		Comment:  *paymentResp.NonceStr,
		PrepayId: *paymentResp.PrepayId,
	}, nil
}
