package payment

import (
	"context"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type RefundPayload = refunddomestic.CreateRequest

type RefundResponse = *refunddomestic.Refund

func Refund(cfg PaymentConfig, payload RefundPayload) (RefundResponse, error) {
	var (
		mchID                      string = cfg.MchId
		mchCertificateSerialNumber string = cfg.MchCertSerialNumber
		mchAPIv3Key                string = cfg.MchApiV3Key
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥 商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client, 并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	svc := refunddomestic.RefundsApiService{Client: client}
	resp, _, err := svc.Create(ctx, payload)

	if err != nil {
		// 处理错误
		return nil, err
	}

	return resp, nil
}
