package payment_test

import (
	"carservice/internal/pkg/wechat/payment"
	"fmt"
	"testing"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

func TestRefund(t *testing.T) {
	c := payment.PaymentConfig{
		MchId:          "1900006891",
		PrivateKeyPath: "P:\\payment\\apiclient_key.pem",
	}

	resp, err := payment.Refund(c, payment.RefundPayload{
		// 子商户的商户号，由微信支付生成并下发。服务商模式下必须传递此参数
		SubMchid: core.String("1900000109"),
		// 原支付交易对应的商户订单号
		OutTradeNo: core.String("1217752501201407033233368018"),
		// 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
		OutRefundNo: core.String("1217752501201407033233368018"),
		// 若商户传入，会在下发给用户的退款消息中体现退款原因
		Reason: core.String("测试退款功能"),
		// 异步接收微信支付退款结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。 如果参数中传了notify_url，则商户平台上配置的回调地址将不会生效，优先回调当前传的这个地址。
		NotifyUrl: core.String("xxxxx"),
		// 订单金额信息
		Amount: &refunddomestic.AmountReq{
			// 退款金额，币种的最小单位，只能为整数，不能超过原订单支付金额。
			Refund: core.Int64(1),
			// 原支付交易的订单总金额，币种的最小单位，只能为整数。
			Total: core.Int64(1),
			// 符合ISO 4217标准的三位字母代码，目前只支持人民币：CNY。
			Currency: core.String("CNY"),
		},
	})
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		t.Fail()
		return
	}
	fmt.Printf("%+v\n", resp)
}
